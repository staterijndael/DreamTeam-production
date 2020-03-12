package controller

import (
	"github.com/gin-gonic/gin"
)

func FNSKeyDevSetter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//app2.GetApp().FNSAPIKey = c.Param("fnsKey")
	}
}
