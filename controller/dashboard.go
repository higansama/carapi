package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard/index", gin.H{
		"title": "Bandros Framework",
	})
}
