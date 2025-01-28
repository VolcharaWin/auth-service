package authorization

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/VolcharaWin/auth-service/custom_errors"
	"github.com/VolcharaWin/auth-service/hashing"
	"github.com/VolcharaWin/auth-service/server"
	"github.com/VolcharaWin/auth-service/token"
	"github.com/VolcharaWin/auth-service/user"
	_ "github.com/mattn/go-sqlite3"
)

// errLoginExists := errors.New("this login already exists")
func Login(srv *server.Server, data server.UserData, db *sql.DB) {
	defer close(data.Done)
	c := data.Context
	log.Println("You are trying to log in.")
	exists, err := user.LoginCheck(db, data.Login)
	if err != nil {
		server.RespondWithError(c, 500, http.StatusText(500))
		log.Println(err)
		return
	}
	if !exists {
		server.RespondWithError(c, 404, http.StatusText(404))
		log.Printf("The user %s does not exist.", data.Login)
		return
	}
	success, err := user.Login(db, data.Login, data.Password)
	if err != nil || !success {
		server.RespondWithError(c, 401, http.StatusText(401))
		log.Println(err)
		return
	}
	log.Println("The login and password match\nCreating the jwt...")
	token, err := token.CreateToken(data.Login)
	if err != nil {
		server.RespondWithError(c, 500, http.StatusText(500))
		log.Println("error while creating a token: ", err)
	} else {
		c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, true)
		server.RespondWithSuccess(c, 200, "successful login")
	}
}

func Registration(srv *server.Server, data server.UserData, db *sql.DB) {
	defer close(data.Done)
	c := data.Context
	log.Println("You are trying to register an account.")
	exists, err := user.LoginCheck(db, data.Login)
	if err != nil {
		server.RespondWithError(c, 500, http.StatusText(500))
		log.Println(err)
	}
	if exists {
		server.RespondWithError(c, 409, custom_errors.ErrLoginExists.Error())
		log.Println("this login already exists")
		return
	}
	hashpass, err := hashing.HashPassword(data.Password)
	if err != nil {
		server.RespondWithError(c, 500, http.StatusText(500))
		log.Println(err)
	}
	success, err := user.Registration(db, data.Login, hashpass)
	if err != nil || !success {
		server.RespondWithError(c, 500, http.StatusText(500))
		log.Println(err)
		return
	}
	server.RespondWithSuccess(c, 201, http.StatusText(201))
	log.Println("Successful registration")
}
