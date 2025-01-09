package server

import (
	"net/http"

	"examples.com/auth-service/token"
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

	protected := r.Group("/protected")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
		})
	}
	return r
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token not found"})
			c.Abort()
			return
		}

		err = token.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func (s *Server) Run() {
	r := s.SetupRouter()
	r.Run("localhost:8082")
}
