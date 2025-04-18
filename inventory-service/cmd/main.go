package main

import (
	"fmt"
	"net"
	"os"

	"E-Commerce/inventory-service/internal/handler"
	"E-Commerce/inventory-service/internal/repository"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	repo := repository.NewProductRepository(db)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryGRPCServer(svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, h)
	fmt.Println("Inventory Service running on :50051")
	s.Serve(lis)
}