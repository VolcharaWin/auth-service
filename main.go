package main

import (
	"errors"
	"log"
	"net/http"

	"examples.com/auth-service/hashing"
	"examples.com/auth-service/server"
	"examples.com/auth-service/storage"
	"examples.com/auth-service/user"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

func main() {
	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading env file: %v\n", err)
	}

	// !!!!!! getting the session key !!!!!!
	//sessionKey := os.Getenv("SESSION_KEY")
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	//Custom errors
	errExistingLogin := errors.New("this login already exists")

	log.Println("====== Initialazing database and connecting to it ======")
	db, err := storage.CreateDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("====== Connection to database is established =======")

	log.Println("====== Setting up the router ======")
	srv := server.NewServer()
	//Обработка входа в систему
	go func() {
		for data := range srv.LoginDataChannel {
			go func(data server.UserData) {
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
				server.RespondWithSuccess(c, 200, "successful login")
				log.Println("Successful authorization")
			}(data)
		}
	}()
	//Обработка регистрации
	go func() {
		for data := range srv.RegisterDataChannel {
			go func(data server.UserData) {
				defer close(data.Done)
				c := data.Context
				log.Println("You are trying to register an account.")
				exists, err := user.LoginCheck(db, data.Login)
				if err != nil {
					server.RespondWithError(c, 500, http.StatusText(500))
					log.Println(err)
				}
				if exists {
					server.RespondWithError(c, 409, http.StatusText(409))
					log.Println(errExistingLogin)
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
			}(data)
		}
	}()
	srv.Run()

}
