package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	c.JSON(http.StatusOK, gin.H{"login": login, "password": password})
}
func SetupRouter() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLFiles("templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/login", login)
	r.Run("localhost:8082")
}
