package main

import (
	"log"
	"time"

	"examples.com/auth-service/check"
	"examples.com/auth-service/server"
	"examples.com/auth-service/storage"
	_ "github.com/golang-jwt/jwt/v5"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

func main() {

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
				log.Printf("login: %s\tpassword: %s\n", data.Login, data.Password)
				exists, err := check.LoginCheck(db, data.Login, data.Password)
				if err != nil {
					log.Println(err)
					return
				}
				if exists {
					log.Printf("The user %s does exist.", data.Login)
				} else {
					log.Printf("The user %s does not exist.", data.Login)
				}
			}(data)
		}
	}()
	//Обработка регистрации
	go func() {
		for data := range srv.RegisterDataChannel {
			go func(data server.UserData) {
				log.Println("You are trying to register an account.")
				time.Sleep(5 * time.Second)
				log.Printf("login: %s\tpassword: %s\n", data.Login, data.Password)
			}(data)
		}
	}()
	srv.Run()

}
