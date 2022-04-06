package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gozknight.com/online-judge/internal/model"
	"gozknight.com/online-judge/internal/util"
	"net/http"
)

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /v1/login [post]
func Login(c *gin.Context) {
	name := c.PostForm("username")
	password := c.PostForm("password")
	// md5
	password = util.MD5(password)
	user := new(model.UserBasic)
	err := model.ORM.Where("name = ? and password = ?", name, password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}
	token, err := util.GenerateToken(user.Identity, user.Name)
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
			"token": token,
		},
	})
}

// GetUser
// @Tags 公共方法
// @Summary 用户详情
// @Accept json
// @Produce json
// @Param identity path string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/{identity} [get]
func GetUser(c *gin.Context) {
	identity := c.Param("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户不存在",
		})
		return
	}
	user := new(model.UserBasic)
	err := model.ORM.Omit("password").Where("identity=?", identity).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
	})
}
