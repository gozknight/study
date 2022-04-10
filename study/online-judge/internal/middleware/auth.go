package middleware

import (
	"github.com/gin-gonic/gin"
	"gozknight.com/online-judge/internal/util"
	"net/http"
)

// AuthAdmin 管理员验证
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("authorization")
		userClaim, err := util.AnalyzeToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized err",
			})
			return
		}
		if userClaim == nil || userClaim.IsAdmin != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized not admin",
			})
			return
		}
		c.Next()
	}
}
