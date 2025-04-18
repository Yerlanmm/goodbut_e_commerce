package service

import (
    "context"
    "github.com/google/uuid"
    pbInventory "E-Commerce/inventory-service/proto"
    "E-Commerce/order-service/internal/entity"
    "E-Commerce/order-service/internal/repository"
    "time"
)

type OrderService interface {
    CreateOrder(userID string, items []*entity.OrderItem) (*entity.Order, error)
    GetOrder(id uuid.UUID) (*entity.Order, error)
    UpdateOrderStatus(id uuid.UUID, status string) error
    ListOrders(userID string, page, pageSize int) ([]*entity.Order, int, error)
}

type orderService struct {
    repo           repository.OrderRepository
    inventoryClient pbInventory.InventoryServiceClient
}

func NewOrderService(repo repository.OrderRepository, inventoryClient pbInventory.InventoryServiceClient) OrderService {
    return &orderService{repo: repo, inventoryClient: inventoryClient}
}

func (s *orderService) CreateOrder(userID string, items []*entity.OrderItem) (*entity.Order, error) {
    for _, item := range items {
        pid := item.ProductID
        p, err := s.inventoryClient.GetProduct(context.Background(), &pbInventory.GetProductRequest{Id: pid.String()})
        if err != nil {
            return nil, err
        }
        item.Price = float64(p.Product.Price)

        resp, err := s.inventoryClient.CheckStock(context.Background(), &pbInventory.CheckStockRequest{
            ProductId: pid.String(),
            Quantity:  int32(item.Quantity),
        })
        if err != nil || !resp.Available {
            return nil, err
        }
    }

    for _, item := range items {
        _, err := s.inventoryClient.UpdateStock(context.Background(), &pbInventory.UpdateStockRequest{
            ProductId: item.ProductID.String(),
            Quantity:  -int32(item.Quantity),
        })
        if err != nil {
            return nil, err
        }
    }

    total := 0.0
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }

    order := &entity.Order{
        ID:          uuid.New(),
        UserID:      userID,
        Status:      "pending",
        TotalAmount: total,
        CreatedAt:   time.Now(),
        Items:       items,
    }

    err := s.repo.Create(order)
    return order, err
}

func (s *orderService) GetOrder(id uuid.UUID) (*entity.Order, error) {
    return s.repo.Get(id)
}

func (s *orderService) UpdateOrderStatus(id uuid.UUID, status string) error {
    return s.repo.UpdateStatus(id, status)
}

func (s *orderService) ListOrders(userID string, page, pageSize int) ([]*entity.Order, int, error) {
    return s.repo.List(userID, page, pageSize)
}