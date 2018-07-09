package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Web struct {}

func (w *Web) Start() {
	router := gin.Default()

	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/img", "ui/content/image")
	router.Static("/css", "ui/content/css")
	router.Static("/js", "ui/content/js")
	router.GET("/", Index)

	router.Run(":8081")
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H {
		"title": "Posts",
	})
}
