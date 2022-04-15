package service

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gozknight.com/online-judge/internal/model"
	"gozknight.com/online-judge/internal/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit/list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", util.DefaultSize))
	if err != nil {
		log.Println(err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", util.DefaultPage))
	if err != nil {
		log.Println(err)
		return
	}
	var count int64
	page = (page - 1) * size
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	list := make([]*model.SubmitBasic, 0)
	tx := model.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count":  count,
			"submit": list,
		},
	})
}

// Submit
// @Tags V1
// @Summary 提交问题
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error" + err.Error(),
		})
		return
	}
	path, err := util.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Save Code Error" + err.Error(),
		})
		return
	}
	u, _ := c.Get("user")
	userClaim := u.(*util.UerClaims)
	sb := &model.SubmitBasic{
		Identity:        util.GetUuid(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}
	pb := new(model.ProblemBasic)
	err = model.ORM.Where("identity = ?", problemIdentity).Preload("TestCase").First(pb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Error" + err.Error(),
		})
		return
	}
	WA := make(chan int)
	OOM := make(chan int)
	CE := make(chan int)
	var msg string
	passCnt := 0
	var mu sync.Mutex
	for _, testC := range pb.TestCase {
		test := testC
		go func() {
			// go run code-user/main.go
			cmd := exec.Command("go", "run", path)
			var out, stdErr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stdErr
			pipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
				return
			}
			io.WriteString(pipe, test.Input)
			var me runtime.MemStats
			runtime.ReadMemStats(&me)
			// 根据输入的测试案例，拿到输入结果和标准输出结果对比
			if err := cmd.Run(); err != nil {
				if err.Error() == "exit status 2" {
					CE <- 1
					msg = stdErr.String()
					return
				}
				return
			}
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			if em.Alloc/1024-me.Alloc/1024 > uint64(pb.MaxMemory) {
				OOM <- 1
				msg = "内存溢出"
				return
			}
			if out.String() != test.Output {
				WA <- 1
				msg = "答案错误"
				return
			}
			mu.Lock()
			defer mu.Unlock()
			passCnt++
		}()
	}
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5
	case <-time.After(time.Millisecond * time.Duration(int64(pb.MaxRuntime))):
		if passCnt == len(pb.TestCase) {
			sb.Status = 1
			msg = "答案正确"
		} else {
			msg = "超出时间限制"
			sb.Status = 3
		}
	}
	if err := model.ORM.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(sb).Error
		if err != nil {
			return err
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		err = tx.Model(new(model.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return err
		}
		m1 := make(map[string]interface{})
		m1["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m1["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		fmt.Println(m)
		fmt.Println(m1)
		err = tx.Model(new(model.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m1).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Error" + err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
	return
}
