package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gozknight.com/online-judge/internal/model"
	"gozknight.com/online-judge/internal/util"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 查看所有问题
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem/list [get]
func GetProblemList(c *gin.Context) {
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
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")
	list := make([]*model.ProblemBasic, 0)
	tx := model.GetProblemList(keyword, categoryIdentity)
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblem
// @Tags 公共方法
// @Summary 问题详情
// @Accept json
// @Produce json
// @Param identity path string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem/{identity} [get]
func GetProblem(c *gin.Context) {
	identity := c.Param("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	data := new(model.ProblemBasic)
	err := model.ORM.Where("identity = ?", identity).Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "当前问题为空",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生未知错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"data": data,
	})
}

// AddProblem
// @Tags 私有方法
// @Summary 添加问题
// @Param authorization header string true "authorization"
// @Param title formData string true "tile"
// @Param content formData string true "content"
// @Param max_runtime formData int true "max_runtime"
// @Param max_memory formData int true "max_memory"
// @Param category_ids formData []string true "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/problem/add [put]
func AddProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(c.PostForm("max_memory"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := util.GetUuid()
	problem := &model.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMemory:  maxMemory,
	}
	categoryBasics := make([]*model.ProblemCategory, 0)
	for _, id := range categoryIds {
		cid, _ := strconv.Atoi(id)
		categoryBasics = append(categoryBasics, &model.ProblemCategory{
			ProblemId:  problem.ID,
			CategoryId: uint(cid),
		})
	}
	problem.ProblemCategories = categoryBasics

	testCaseBasics := make([]*model.TestCase, 0)
	for _, test := range testCases {
		testCaseMap := make(map[string]string)
		err := json.Unmarshal([]byte(test), &testCaseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}
		if _, ok := testCaseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}
		if _, ok := testCaseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}
		testCaseBasics = append(testCaseBasics, &model.TestCase{
			Identity:        util.GetUuid(),
			ProblemIdentity: identity,
			Input:           testCaseMap["input"],
			Output:          testCaseMap["output"],
		})
	}
	problem.TestCase = testCaseBasics
	err := model.ORM.Create(problem).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建问题失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题创建成功",
	})

}

// EditProblem
// @Tags 私有方法
// @Summary 修改问题
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int true "max_runtime"
// @Param max_memory formData int false "max_memory"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/problem/edit [post]
func EditProblem(c *gin.Context) {
	identity := c.Query("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(c.PostForm("max_memory"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := model.ORM.Transaction(func(tx *gorm.DB) error {
		problemBasic := &model.ProblemBasic{
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMemory:  maxMemory,
		}
		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		err = tx.Where("identity = ?", identity).Find(problemBasic).Error
		if err != nil {
			return err
		}
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(model.ProblemCategory)).Error
		if err != nil {
			return err
		}
		var pcs []*model.ProblemCategory
		for _, id := range categoryIds {
			intId, _ := strconv.Atoi(id)
			pcs = append(pcs, &model.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		err = tx.Where("problem_identity = ?", identity).Delete(new(model.TestCase)).Error
		if err != nil {
			return err
		}
		var tcs []*model.TestCase
		for _, testcase := range testCases {
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testcase), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例输入格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例输出格式错误")
			}
			tcs = append(tcs, &model.TestCase{
				Identity:        util.GetUuid(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改问题失败，" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "修改问题成功，",
	})
	return
}
