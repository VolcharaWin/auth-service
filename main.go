package main

import (
	"errors"
	"log"

	"examples.com/auth-service/hashing"
	"examples.com/auth-service/server"
	"examples.com/auth-service/storage"
	"examples.com/auth-service/user"
	_ "github.com/golang-jwt/jwt/v5"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

func main() {
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
				log.Println("You are trying to log in.")
				exists, err := user.LoginCheck(db, data.Login)
				if err != nil {
					log.Println(err)
					return
				}
				if !exists {
					log.Printf("The user %s does not exist.", data.Login)
					return
				}
				success, err := user.Login(db, data.Login, data.Password)
				if err != nil || !success {
					log.Println(err)
					return
				}
				log.Println("Successful authorization")
			}(data)
		}
	}()
	//Обработка регистрации
	go func() {
		for data := range srv.RegisterDataChannel {
			go func(data server.UserData) {
				log.Println("You are trying to register an account.")
				exists, err := user.LoginCheck(db, data.Login)
				if err != nil {
					log.Println(err)
				}
				if exists {
					log.Println(errExistingLogin)
					return
				}
				hashpass, err := hashing.HashPassword(data.Password)
				if err != nil {
					log.Println(err)
				}
				success, err := user.Registration(db, data.Login, hashpass)
				if err != nil || !success {
					log.Println(err)
					return
				}
				log.Println("Successful registration")
			}(data)
		}
	}()
	srv.Run()

}
