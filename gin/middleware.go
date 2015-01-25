package gin

import (
	"github.com/gin-gonic/gin"
	"yaag/middleware"
)

func Yaag() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
