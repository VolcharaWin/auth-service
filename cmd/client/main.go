package main

import (
	"context"
	"log"

	pb "github.com/VolcharaWin/auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	registerResp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Email:    "user@example.com",
		Password: "securepassword",
	})

	if err != nil {
		log.Fatalf("registration failed: %v", err)
	}
	log.Printf("Registered user ID: %s\n", registerResp.UserId)

	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
		Email:    "user@example.com",
		Password: "securepassword",
	})

	if err != nil {
		log.Fatalf("login failed: %v", err)
	}

	log.Printf("Auth token: %s\n", loginResp.Token)

	validateResp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: loginResp.Token,
	})

	if err != nil {
		log.Fatalf("token validation failed: %v", err)
	}
	log.Printf("Token vlaid: %t, User ID: %s\n", validateResp.Valid, validateResp.UserId)
}
