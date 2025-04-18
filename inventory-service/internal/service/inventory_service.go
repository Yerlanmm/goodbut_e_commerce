package service

import (
    "github.com/google/uuid"
    "E-Commerce/inventory-service/internal/entity"
    "E-Commerce/inventory-service/internal/repository"
)

type InventoryService interface {
    CreateProduct(p *entity.Product) error
    GetProduct(id uuid.UUID) (*entity.Product, error)
    UpdateProduct(p *entity.Product) error
    DeleteProduct(id uuid.UUID) error
    ListProducts(categoryID string, page, pageSize int) ([]*entity.Product, int, error)
    CheckStock(productID uuid.UUID, quantity int) (bool, error)
    UpdateStock(productID uuid.UUID, quantity int) error
}

type inventoryService struct {
    repo repository.ProductRepository
}

func NewInventoryService(repo repository.ProductRepository) InventoryService {
    return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateProduct(p *entity.Product) error {
    p.ID = uuid.New()
    return s.repo.Create(p)
}

func (s *inventoryService) GetProduct(id uuid.UUID) (*entity.Product, error) {
    return s.repo.Get(id)
}

func (s *inventoryService) UpdateProduct(p *entity.Product) error {
    return s.repo.Update(p)
}

func (s *inventoryService) DeleteProduct(id uuid.UUID) error {
    return s.repo.Delete(id)
}

func (s *inventoryService) ListProducts(categoryID string, page, pageSize int) ([]*entity.Product, int, error) {
    return s.repo.List(categoryID, page, pageSize)
}

func (s *inventoryService) CheckStock(productID uuid.UUID, quantity int) (bool, error) {
    return s.repo.CheckStock(productID, quantity)
}

func (s *inventoryService) UpdateStock(productID uuid.UUID, quantity int) error {
    return s.repo.UpdateStock(productID, quantity)
}