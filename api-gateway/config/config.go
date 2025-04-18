package config

import (
	pbInventory "E-Commerce/inventory-service/proto"
	pbOrder "E-Commerce/order-service/proto"
	"google.golang.org/grpc"
)

type Config struct {
	InventoryClient pbInventory.InventoryServiceClient
	OrderClient     pbOrder.OrderServiceClient
}

func NewConfig() *Config {
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &Config{
		InventoryClient: pbInventory.NewInventoryServiceClient(inventoryConn),
		OrderClient:     pbOrder.NewOrderServiceClient(orderConn),
	}
}