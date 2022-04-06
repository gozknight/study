package service

import (
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
// @Router /problem [get]
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
