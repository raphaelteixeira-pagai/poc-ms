package httphdl

import (
	"github.com/gin-gonic/gin"
	"github.com/golangsugar/chatty"
)

func MetricMiddleware(escopo string, tipo string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			go chatty.Infof("scope [%v] type [%v]", escopo, tipo)
		}()
		c.Next()
	}
}
