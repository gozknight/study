package service

import (
	"encoding/json"
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
// @Param max_runtime formData string true "max_runtime"
// @Param max_memory formData string true "max_memory"
// @Param category_ids formData array true "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/problem/add [put]
func AddProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime := c.PostForm("max_runtime")
	maxMemory := c.PostForm("max_memory")
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || maxRuntime == "" || maxMemory == "" || len(categoryIds) == 0 || len(testCases) == 0 {
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
