package sessions

import (
	"net/http"

	"examples.com/auth-service/server"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func NewSession(c *gin.Context, r *gin.Engine, login, secret_key string) {
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(secret_key))
	if err != nil {
		server.RespondWithError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	r.Use(sessions.Sessions("my-session", store))

}
