package main

import (
	"log"
	"time"

	"examples.com/auth-service/server"
	"examples.com/auth-service/storage"
	_ "github.com/golang-jwt/jwt/v5"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
)

func main() {
	log.Println("====== Initialazing database and connecting to it ======")
	_, err := storage.CreateDB()
	if err != nil {
		panic(err)
	}
	//defer con.Close()
	log.Println("====== Connection to database is established =======")
	log.Println("====== Setting up the router ======")
	srv := server.NewServer()
	//Обработка входа в систему
	go func() {
		for data := range srv.LoginDataChannel {
			log.Println("You are trying to log in.")
			time.Sleep(5 * time.Second)
			log.Printf("login: %s\tpassword: %s\n", data.Login, data.Password)
		}
	}()
	//Обработка регистрации
	go func() {
		for data := range srv.RegisterDataChannel {
			log.Println("You are trying to register an account.")
			time.Sleep(5 * time.Second)
			log.Printf("login: %s\tpassword: %s\n", data.Login, data.Password)
		}
	}()
	srv.Run()

}
