// product-service/main.go (partial)
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/gorm"
	
	pb "github.com/yourusername/myapp/product-service/proto"
	"github.com/yourusername/myapp/product-service/models"
)

type productServer struct {
	pb.UnimplementedProductServiceServer
	db *gorm.DB
}

// Implementasi CreateProduct
func (s *productServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	// Create product in database
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedByID: req.UserId,
	}
	
	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}
	
	return &pb.Product{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedBy:   product.CreatedByID.String(),
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}, nil
}

// Implementasi method lainnya serupa