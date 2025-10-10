package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello 首页，用于健康检查
func Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello, word.")
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
