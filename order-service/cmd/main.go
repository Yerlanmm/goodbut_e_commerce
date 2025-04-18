package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pbInventory "E-Commerce/inventory-service/proto"
	"E-Commerce/order-service/internal/handler"
	"E-Commerce/order-service/internal/repository"
	"E-Commerce/order-service/internal/service"
	pb "E-Commerce/order-service/proto"
	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	ic := pbInventory.NewInventoryServiceClient(conn)
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, ic)
	h := handler.NewOrderGRPCServer(svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, h)
	fmt.Println("Order Service running on :50052")
	s.Serve(lis)
}
