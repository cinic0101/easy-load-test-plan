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
	router.Static("/component", "ui/content/components")

	router.GET("/", Index)
	router.GET("/prototype", Prototype)

	router.POST("/service/project", AddNewProject)

	router.Run(":8081")
}

func Index(c *gin.Context) {
	storage.Load()

	c.HTML(http.StatusOK, "index.html", gin.H {
		"projects": storage.AllProjects.Projects,
	})
}

func Prototype(c *gin.Context) {
	c.HTML(http.StatusOK, "prototype.html", gin.H {})
}

func AddNewProject(c *gin.Context)  {
	id := c.PostForm("id")
	name := c.PostForm("name")
	desc := c.PostForm("desc")

	storage.AddNewTestPlan(id, name, desc)

	c.String(http.StatusOK, "%v", true)
}
