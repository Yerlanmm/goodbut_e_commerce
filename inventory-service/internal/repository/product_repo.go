package repository

import (
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "E-Commerce/inventory-service/internal/entity"
)

type ProductRepository interface {
    Create(p *entity.Product) error
    Get(id uuid.UUID) (*entity.Product, error)
    Update(p *entity.Product) error
    Delete(id uuid.UUID) error
    List(categoryID string, page, pageSize int) ([]*entity.Product, int, error)
    CheckStock(productID uuid.UUID, quantity int) (bool, error)
    UpdateStock(productID uuid.UUID, quantity int) error
}

type productRepository struct {
    db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(product *entity.Product) error {
    query := `INSERT INTO products (id, name, description, price, stock, category_id) 
              VALUES (:id, :name, :description, :price, :stock, :category_id)`
    _, err := r.db.NamedExec(query, product)
    return err
}

func (r *productRepository) Get(id uuid.UUID) (*entity.Product, error) {
    var p entity.Product
    err := r.db.Get(&p, "SELECT * FROM products WHERE id = $1", id)
    if err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *productRepository) Update(product *entity.Product) error {
    query := `UPDATE products SET name = :name, description = :description, price = :price, 
              stock = :stock, category_id = :category_id WHERE id = :id`
    _, err := r.db.NamedExec(query, product)
    return err
}

func (r *productRepository) Delete(id uuid.UUID) error {
    _, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
    return err
}

func (r *productRepository) List(categoryID string, page, pageSize int) ([]*entity.Product, int, error) {
    var products []*entity.Product
    var total int
    query := "SELECT * FROM products"
    countQuery := "SELECT COUNT(*) FROM products"
    if categoryID != "" {
        query += " WHERE category_id = $1"
        countQuery += " WHERE category_id = $1"
        err := r.db.Get(&total, countQuery, categoryID)
        if err != nil {
            return nil, 0, err
        }
        err = r.db.Select(&products, query+" LIMIT $2 OFFSET $3", categoryID, pageSize, (page-1)*pageSize)
        if err != nil {
            return nil, 0, err
        }
    } else {
        err := r.db.Get(&total, countQuery)
        if err != nil {
            return nil, 0, err
        }
        err = r.db.Select(&products, query+" LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
        if err != nil {
            return nil, 0, err
        }
    }
    return products, total, nil
}

func (r *productRepository) CheckStock(productID uuid.UUID, quantity int) (bool, error) {
    var stock int
    err := r.db.Get(&stock, "SELECT stock FROM products WHERE id = $1", productID)
    if err != nil {
        return false, err
    }
    return stock >= quantity, nil
}

func (r *productRepository) UpdateStock(productID uuid.UUID, quantity int) error {
    _, err := r.db.Exec("UPDATE products SET stock = stock + $1 WHERE id = $2", quantity, productID)
    return err
}