// auth-service/main.go (partial)
package main

import (
	"context"
	"log"
	"net"
	
	"google.golang.org/grpc"
	"firebase.google.com/go/v4/auth"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	
	pb "github.com/raflibima25/go-belajar-kong-firebase-grpc/auth-service/proto"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	firebaseAuth *auth.Client
	db           *gorm.DB
}

// Implementasi Register
func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	// 1. Create user in Firebase Auth
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.Name)
	
	firebaseUser, err := s.firebaseAuth.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	
	// 2. Create user in our database
	user := models.User{
		FirebaseUID: firebaseUser.UID,
		Email:       req.Email,
		Role:        "user", // Default role is user
	}
	
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	
	// 3. Generate custom token
	token, err := s.firebaseAuth.CustomToken(ctx, firebaseUser.UID)
	if err != nil {
		return nil, err
	}
	
	return &pb.AuthResponse{
		Token:  token,
		UserId: user.ID.String(),
		Role:   user.Role,
	}, nil
}

// Implementasi Login dan Validate serupa