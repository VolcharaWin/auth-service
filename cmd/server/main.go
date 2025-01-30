package main

import (
	"context"
	"crypto/rand"
	"net"
	"time"

	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/VolcharaWin/auth-service/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	users         map[string]*User
	jwtSecret     []byte
	tokenDuration time.Duration
}

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewAuthServer(secret []byte, tokenDuration time.Duration) *AuthServer {
	return &AuthServer{
		users:         make(map[string]*User),
		jwtSecret:     secret,
		tokenDuration: tokenDuration,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if _, exists := s.users[req.Email]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	userID, _ := generateRandomString(16)
	s.users[req.Email] = &User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	return &pb.RegisterResponse{UserId: userID}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, exists := s.users[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate code")
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := s.validateJWT(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
	}, nil
}

func (s *AuthServer) validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Internal, "unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "invalid token")
}
func (s *AuthServer) generateJWT(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func main() {
	jwtSecret := []byte(os.Getenv("SESSION_KEY"))
	tokenDuration := 24 * time.Hour

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	authServer := NewAuthServer(jwtSecret, tokenDuration)
	pb.RegisterAuthServiceServer(server, authServer)

	log.Printf("Server starting on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
