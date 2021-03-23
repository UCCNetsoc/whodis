package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitAPI() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.Run(":8080")
}
