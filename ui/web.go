package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Web struct {}

var storage *Storage

func (w *Web) Start() {
	storage = new(Storage)
	storage.CreateIfNotExists()

	router := gin.Default()
	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/img", "ui/content/image")
	router.Static("/css", "ui/content/css")
	router.Static("/js", "ui/content/js")
	router.GET("/", Index)
	router.GET("/prototype", Prototype)

	router.Run(":8081")
}

func Index(c *gin.Context) {
	storage.Load()
	storage.AddNewTestPlan("1234", "abcd")

	c.HTML(http.StatusOK, "index.html", gin.H {})
}

func Prototype(c *gin.Context) {
	c.HTML(http.StatusOK, "prototype.html", gin.H {})
}

