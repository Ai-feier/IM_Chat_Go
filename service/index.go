package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetIndex
// @Tags 首页
// @Accept json
// @Success 200 {string} Welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
