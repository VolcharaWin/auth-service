package main

import (
	"log"

	"github.com/VolcharaWin/auth-service/authorization"
	"github.com/VolcharaWin/auth-service/server"
	"github.com/VolcharaWin/auth-service/storage"
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
				authorization.Login(srv, data, db)
			}(data)
		}
	}()
	//Обработка регистрации
	go func() {
		for data := range srv.RegisterDataChannel {
			go func(data server.UserData) {
				authorization.Registration(srv, data, db)
			}(data)
		}
	}()
	srv.Run()

}
