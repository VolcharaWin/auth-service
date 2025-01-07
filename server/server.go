package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	Login    string
	Password string
	Context  *gin.Context
	Done     chan bool
}

type Server struct {
	LoginDataChannel    chan UserData
	RegisterDataChannel chan UserData
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func RespondWithSuccess(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func (s *Server) register(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	done := make(chan bool)

	s.RegisterDataChannel <- UserData{
		Login:    login,
		Password: password,
		Context:  c,
		Done:     done,
	}

	<-done
}
func (s *Server) login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	done := make(chan bool)

	s.LoginDataChannel <- UserData{
		Login:    login,
		Password: password,
		Context:  c,
		Done:     done,
	}

	<-done
}
func NewServer() *Server {
	return &Server{
		LoginDataChannel:    make(chan UserData),
		RegisterDataChannel: make(chan UserData),
	}
}
func (s *Server) SetupRouter() *gin.Engine {
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
	r.POST("/login", s.login)
	r.POST("/register", s.register)
	return r
}

func (s *Server) Run() {
	r := s.SetupRouter()
	r.Run("localhost:8082")
}
