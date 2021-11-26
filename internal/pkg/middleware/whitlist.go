package middleware

import (
	"net/http"
	"wxrobot/internal/pkg/access"

	"github.com/gin-gonic/gin"
)

func WhiteListCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _ := c.RemoteIP()
		isWhite, err := access.IsCallbackWhiteList(ip.String())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		if !isWhite {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}
