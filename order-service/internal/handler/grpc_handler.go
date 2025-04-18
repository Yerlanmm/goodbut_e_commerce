package handler

import (
    "context"
    "github.com/google/uuid"
    pb "E-Commerce/order-service/proto"
    "E-Commerce/order-service/internal/service"
    "E-Commerce/order-service/internal/entity"
)

type OrderGRPCServer struct {
    pb.UnimplementedOrderServiceServer
    svc service.OrderService
}

func NewOrderGRPCServer(svc service.OrderService) *OrderGRPCServer {
    return &OrderGRPCServer{svc: svc}
}

func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
    items := make([]*entity.OrderItem, len(req.Items))
    for i, item := range req.Items {
        pid, err := uuid.Parse(item.ProductId)
        if err != nil {
            return nil, err
        }
        items[i] = &entity.OrderItem{
            ProductID: pid,
            Quantity:  int(item.Quantity),
        }
    }
    order, err := s.svc.CreateOrder(req.UserId, items)
    if err != nil {
        return nil, err
    }
    pbOrder := &pb.Order{
        Id:          order.ID.String(),
        UserId:      order.UserID,
        Status:      order.Status,
        TotalAmount: float32(order.TotalAmount),
        Items:       make([]*pb.OrderItem, len(order.Items)),
    }
    for i, item := range order.Items {
        pbOrder.Items[i] = &pb.OrderItem{
            ProductId: item.ProductID.String(),
            Quantity:  int32(item.Quantity),
            Price:     float32(item.Price),
        }
    }
    return &pb.CreateOrderResponse{Order: pbOrder}, nil
}

func (s *OrderGRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
    id, err := uuid.Parse(req.Id)
    if err != nil {
        return nil, err
    }
    order, err := s.svc.GetOrder(id)
    if err != nil {
        return nil, err
    }
    pbOrder := &pb.Order{
        Id:          order.ID.String(),
        UserId:      order.UserID,
        Status:      order.Status,
        TotalAmount: float32(order.TotalAmount),
        Items:       make([]*pb.OrderItem, len(order.Items)),
    }
    for i, item := range order.Items {
        pbOrder.Items[i] = &pb.OrderItem{
            ProductId: item.ProductID.String(),
            Quantity:  int32(item.Quantity),
            Price:     float32(item.Price),
        }
    }
    return &pb.GetOrderResponse{Order: pbOrder}, nil
}

func (s *OrderGRPCServer) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
    id, err := uuid.Parse(req.Id)
    if err != nil {
        return nil, err
    }
    err = s.svc.UpdateOrderStatus(id, req.Status)
    if err != nil {
        return nil, err
    }
    order, err := s.svc.GetOrder(id)
    if err != nil {
        return nil, err
    }
    pbOrder := &pb.Order{
        Id:          order.ID.String(),
        UserId:      order.UserID,
        Status:      order.Status,
        TotalAmount: float32(order.TotalAmount),
        Items:       make([]*pb.OrderItem, len(order.Items)),
    }
    for i, item := range order.Items {
        pbOrder.Items[i] = &pb.OrderItem{
            ProductId: item.ProductID.String(),
            Quantity:  int32(item.Quantity),
            Price:     float32(item.Price),
        }
    }
    return &pb.UpdateOrderResponse{Order: pbOrder}, nil
}

func (s *OrderGRPCServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    orders, total, err := s.svc.ListOrders(req.UserId, int(req.Page), int(req.PageSize))
    if err != nil {
        return nil, err
    }
    pbOrders := make([]*pb.Order, len(orders))
    for i, order := range orders {
        pbOrders[i] = &pb.Order{
            Id:          order.ID.String(),
            UserId:      order.UserID,
            Status:      order.Status,
            TotalAmount: float32(order.TotalAmount),
            Items:       make([]*pb.OrderItem, len(order.Items)),
        }
        for j, item := range order.Items {
            pbOrders[i].Items[j] = &pb.OrderItem{
                ProductId: item.ProductID.String(),
                Quantity:  int32(item.Quantity),
                Price:     float32(item.Price),
            }
        }
    }
    return &pb.ListOrdersResponse{
        Orders: pbOrders,
        Total:  int32(total),
    }, nil
}