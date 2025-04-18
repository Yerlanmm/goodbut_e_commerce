package handler

import (
	"E-Commerce/inventory-service/internal/entity"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"
	"context"

	"github.com/google/uuid"
)

type InventoryGRPCServer struct {
	pb.UnimplementedInventoryServiceServer
	svc service.InventoryService
}

func NewInventoryGRPCServer(svc service.InventoryService) *InventoryGRPCServer {
	return &InventoryGRPCServer{svc: svc}
}

func (s *InventoryGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  req.CategoryId, 
	}

	err := s.svc.CreateProduct(p)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	p, err := s.svc.GetProduct(id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	p := &entity.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  req.CategoryId, 
	}

	err = s.svc.UpdateProduct(p)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	err = s.svc.DeleteProduct(id)
	return &pb.DeleteProductResponse{Success: err == nil}, err
}

func (s *InventoryGRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.svc.ListProducts(req.CategoryId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, p := range products {
		pbProducts[i] = &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		}
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}, nil
}

func (s *InventoryGRPCServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	pid, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, err
	}
	available, err := s.svc.CheckStock(pid, int(req.Quantity))
	if err != nil {
		return nil, err
	}
	return &pb.CheckStockResponse{Available: available}, nil
}

func (s *InventoryGRPCServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	pid, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, err
	}
	err = s.svc.UpdateStock(pid, int(req.Quantity))
	return &pb.UpdateStockResponse{Success: err == nil}, err
}
