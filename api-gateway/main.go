// api-gateway/main.go
package main

import (
	"github.com/gin-gonic/gin"

	authpb "github.com/yourusername/myapp/auth-service/proto"
	productpb "github.com/yourusername/myapp/product-service/proto"
	"google.golang.org/grpc"
)

func main() {
	r := gin.Default()

	// Setup gRPC client connections
	authConn, err := grpc.Dial("auth-service:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)

	productConn, err := grpc.Dial("product-service:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer productConn.Close()
	productClient := productpb.NewProductServiceClient(productConn)

	// Middleware untuk autentikasi
	authMiddleware := func(requiredRole string) gin.HandlerFunc {
		return func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
				return
			}

			// Validate token with auth service
			resp, err := authClient.Validate(c, &authpb.ValidateRequest{
				Token: token,
			})
			if err != nil || !resp.Valid {
				c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
				return
			}

			// Check role if required
			if requiredRole != "" && resp.Role != requiredRole {
				c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
				return
			}

			// Add user info to context
			c.Set("userId", resp.UserId)
			c.Set("role", resp.Role)
			c.Next()
		}
	}

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			var req authpb.RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			resp, err := authClient.Register(c, &req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, resp)
		})

		auth.POST("/login", func(c *gin.Context) {
			// Similar implementation
		})
	}

	// Product routes
	products := r.Group("/products")
	products.Use(authMiddleware("")) // Auth required for all product routes
	{
		products.POST("/", func(c *gin.Context) {
			var req productpb.CreateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Set user ID from auth middleware
			req.UserId = c.GetString("userId")

			resp, err := productClient.CreateProduct(c, &req)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, resp)
		})

		// Other product routes similar
	}

	// Admin-only routes
	admin := r.Group("/admin")
	admin.Use(authMiddleware("admin"))
	{
		// Admin-specific endpoints
	}

	r.Run(":8080")
}
