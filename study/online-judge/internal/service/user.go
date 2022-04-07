package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gozknight.com/online-judge/internal/model"
	"gozknight.com/online-judge/internal/util"
	"log"
	"net/http"
	"strconv"
	"time"
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
	token, err := util.GenerateToken(user.Identity, user.Name, user.IsAdmin)
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

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /v1/send [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	code := util.GetRandomCode()
	err := util.SendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生错误",
		})
		return
	}
	model.RDB.Set(c, email, code, time.Second*300)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"email": email,
			"code":  code,
		},
	})
}

// Register
// @Tags 公共方法
// @Summary 注册用户
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string true "phone"
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /v1/register [post]
func Register(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	code := c.PostForm("code")
	var cnt int64
	err := model.ORM.Where("email = ?", email).Model(new(model.UserBasic)).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生未知错误",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱已被注册",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := model.RDB.Get(c, email).Result()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取验证码失败，请重新获取",
		})
		return
	}
	if sysCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确，请重新输入验证码",
		})
		return
	}
	// 数据插入
	userIdentity := util.GetUuid()
	user := &model.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: util.MD5(password),
		Phone:    phone,
		Email:    email,
	}
	err = model.ORM.Create(user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}
	token, err := util.GenerateToken(userIdentity, name, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成token错误，" + err.Error(),
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

// GetRankList
// @Tags 公共方法
// @Summary 查看排名
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank [get]
func GetRankList(c *gin.Context) {
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
	list := make([]*model.UserBasic, 0)
	err = model.ORM.Model(new(model.UserBasic)).Count(&count).Order("finish_problem_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取排名失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"list":  list,
		},
	})
}
