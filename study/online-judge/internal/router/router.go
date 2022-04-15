package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gozknight.com/online-judge/docs"
	"gozknight.com/online-judge/internal/middleware"
	"gozknight.com/online-judge/internal/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	// 路由规则
	v1 := r.Group("v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		v1.GET("/ping", service.Ping)
	}
	// 问题
	{
		r.GET("/problem/list", service.GetProblemList)
		r.GET("/problem/:identity", service.GetProblem)
	}
	// 排行榜
	{
		r.GET("/rank", service.GetRankList)
	}

	// 用户
	{
		r.GET("/user/:identity", service.GetUser)
		r.POST("/v1/login", service.Login)
		r.POST("/v1/send", service.SendCode)
		r.POST("/v1/register", service.Register)
	}

	// 提交记录
	{
		r.GET("/submit/list", service.GetSubmitList)
	}

	// 管理员
	{
		admin := r.Group("/admin", middleware.AuthAdmin())
		admin.POST("/problem/add", service.AddProblem)
		admin.GET("/category/list", service.GetCategoryList)
		admin.POST("/category/add", service.AddCategory)
		admin.PUT("/category/edit", service.EditCategory)
		admin.DELETE("/category/delete", service.DeleteCategory)
		admin.PUT("/problem/edit", service.EditProblem)
	}

	{
		user := r.Group("/user", middleware.AuthUser())
		user.POST("/submit", service.Submit)
	}
	return r
}
