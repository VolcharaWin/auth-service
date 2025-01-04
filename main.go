package main

import (
	"log"

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
	server.SetupRouter()
}
