package service

import (
	"github.com/gin-gonic/gin"
	"gozknight.com/online-judge/internal/model"
	"gozknight.com/online-judge/internal/util"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 私有方法
// @Summary 查看分类列表
// @Accept json
// @Produce json
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category/list [get]
func GetCategoryList(c *gin.Context) {
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
	var categories *[]model.CategoryBasic
	err = model.ORM.Model(new(model.CategoryBasic)).Where("name like ?", "%"+keyword+"%").Count(&count).
		Offset(page).Limit(size).Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"list":  categories,
		},
	})
	return
}

// AddCategory
// @Tags 私有方法
// @Summary 添加分类
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/category/add [put]
func AddCategory(c *gin.Context) {
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := &model.CategoryBasic{
		Identity: util.GetUuid(),
		Name:     name,
		ParentId: parentId,
	}
	err := model.ORM.Create(category).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分类失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "创建分类成功",
	})
	return
}

// EditCategory
// @Tags 私有方法
// @Summary 修改分类
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Param name formData string true "name"
// @Param parent_id formData string false "parent_id"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/category/edit [post]
func EditCategory(c *gin.Context) {
	identity := c.Query("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := &model.CategoryBasic{
		Name:     name,
		ParentId: parentId,
	}
	var cnt int64
	err := model.ORM.Model(new(model.CategoryBasic)).Where("identity = ?", identity).Count(&cnt).Updates(category).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改分类失败" + err.Error(),
		})
		return
	}
	if cnt <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改分类失败，问题分类不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改分类成功",
	})
	return
}

// DeleteCategory
// @Tags 私有方法
// @Summary 删除分类
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /admin/category/delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	var count int64
	err := model.ORM.Model(new(model.ProblemCategory)).
		Where("category_id = (select id from category_basic where identity = ? limit 1)", identity).
		Count(&count).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败" + err.Error(),
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下面存在问题，无法删除",
		})
		return
	}
	err = model.ORM.Where("identity = ?", identity).Delete(new(model.CategoryBasic)).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  "删除成功",
	})
	return
}
