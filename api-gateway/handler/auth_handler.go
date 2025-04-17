package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/raflibima25/go-belajar-kong-firebase-grpc/grpc/pb"
)

type AuthHandler struct {
	AuthClient pb.AuthServiceClient
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(authClient pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		AuthClient: authClient,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	resp, err := h.AuthClient.Re

}
