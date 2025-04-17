package handler

import (
	"auth-service/model"
	"context"
	"errors"
	"time"

	"firebase.google.com/go/v4/auth"
	pb "github.com/raflibima25/go-belajar-kong-firebase-grpc/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	DB           *gorm.DB
	FirebaseAuth *auth.Client
}

func NewAuthHandler(db *gorm.DB, firebaseAuth *auth.Client) *AuthHandler {
	return &AuthHandler{
		DB:           db,
		FirebaseAuth: firebaseAuth,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	// 1. Create user in Firebase Auth
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password).
		DisplayName(req.Name)

	firebaseUser, err := h.FirebaseAuth.CreateUser(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create Firebase user: %v", err)
	}

	// Set custom claims (role)
	claims := map[string]interface{}{
		"role": "user",
	}
	if err := h.FirebaseAuth.SetCustomUserClaims(ctx, firebaseUser.UID, claims); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to set custom claims: %v", err)
	}

	// 2. Create user in postgres database
	user := model.User{
		FirebaseUID: firebaseUser.UID,
		Email:       req.Email,
		Name:        req.Name,
		Role:        "user",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		// Rollback Firebase user creation if DB insert fails
		h.FirebaseAuth.DeleteUser(ctx, firebaseUser.UID)
		return nil, status.Errorf(codes.Internal, "Failed to create user in database: %v", err)
	}

	// 3. Generate custom token
	token, err := h.FirebaseAuth.CustomTokenWithClaims(ctx, firebaseUser.UID, claims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token: %v", err)
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.User{
			Id:          user.ID.String(),
			FirebaseUid: user.FirebaseUID,
			Email:       user.Email,
			Name:        user.Name,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	userRecord, err := h.FirebaseAuth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found: %v", err)
	}

	var user model.User
	if err := h.DB.Where("firebase_uid = ?", userRecord.UID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found in database: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	claims := map[string]interface{}{
		"role": user.Role,
	}
	token, err := h.FirebaseAuth.CustomTokenWithClaims(ctx, userRecord.UID, claims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token: %v", err)
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.User{
			Id:          user.ID.String(),
			FirebaseUid: user.FirebaseUID,
			Email:       user.Email,
			Name:        user.Name,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	token, err := h.FirebaseAuth.VerifyIDToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Valid: false,
		}, nil
	}

	var user model.User
	if err := h.DB.Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found in database")
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	return &pb.ValidateResponse{
		Valid: true,
		User: &pb.User{
			Id:          user.ID.String(),
			FirebaseUid: user.FirebaseUID,
			Email:       user.Email,
			Name:        user.Name,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *AuthHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.User, error) {
	var user model.User
	if err := h.DB.Where("id = ?", req.Id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	return &pb.User{
		Id:          user.ID.String(),
		FirebaseUid: user.FirebaseUID,
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil

}
