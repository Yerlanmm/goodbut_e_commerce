package main

import (
	"log"

	"E-Commerce/api-gateway/internal/handler"
	"E-Commerce/api-gateway/internal/middleware"
	pbInventory "E-Commerce/inventory-service/proto"
	pbOrder "E-Commerce/order-service/proto"
	pbUser "E-Commerce/user-service/proto"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	} else {
		log.Println("Connected to Inventory Service")
	}
	defer inventoryConn.Close()

	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	} else {
		log.Println("Connected to Order Service")
	}
	defer orderConn.Close()

	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	} else {
		log.Println("Connected to User Service")
	}
	defer userConn.Close()

	inventoryClient := pbInventory.NewInventoryServiceClient(inventoryConn)
	orderClient := pbOrder.NewOrderServiceClient(orderConn)
	userClient := pbUser.NewUserServiceClient(userConn)

	r := gin.Default()

	userHandler := handler.NewUserHandler(userClient)
	h := handler.NewRESTHandler(inventoryClient, orderClient)

	r.POST("/auth/register", userHandler.Register)
	r.POST("/auth/login", userHandler.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		// Admin only routes
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/products", h.CreateProduct)
			admin.PATCH("/products/:id", h.UpdateProduct)
			admin.DELETE("/products/:id", h.DeleteProduct)
			admin.PATCH("/orders/:id", h.UpdateOrder)
		}

		// User routes
		protected.GET("/products/:id", h.GetProduct)
		protected.GET("/products", h.ListProducts)
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)
		protected.GET("/orders", h.ListOrders)

		// Profile routes
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me", userHandler.UpdateCurrentUser)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}