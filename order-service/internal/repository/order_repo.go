package repository

import (
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "E-Commerce/order-service/internal/entity"
)

type orderRepository struct {
    db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *entity.Order) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `INSERT INTO orders (id, user_id, status, total_amount, created_at) 
              VALUES (:id, :user_id, :status, :total_amount, :created_at)`
    _, err = tx.NamedExec(query, order)
    if err != nil {
        return err
    }

    for _, item := range order.Items {
        item.OrderID = order.ID
        query = `INSERT INTO order_items (order_id, product_id, quantity, price) 
                 VALUES (:order_id, :product_id, :quantity, :price)`
        _, err = tx.NamedExec(query, item)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *orderRepository) Get(id uuid.UUID) (*entity.Order, error) {
    var o entity.Order
    err := r.db.Get(&o, "SELECT * FROM orders WHERE id = $1", id)
    if err != nil {
        return nil, err
    }
    err = r.db.Select(&o.Items, "SELECT * FROM order_items WHERE order_id = $1", id)
    if err != nil {
        return nil, err
    }
    return &o, nil
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status string) error {
    _, err := r.db.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, id)
    return err
}

func (r *orderRepository) List(userID string, page, pageSize int) ([]*entity.Order, int, error) {
    var orders []*entity.Order
    var total int
    query := "SELECT * FROM orders WHERE user_id = $1 LIMIT $2 OFFSET $3"
    countQuery := "SELECT COUNT(*) FROM orders WHERE user_id = $1"
    err := r.db.Get(&total, countQuery, userID)
    if err != nil {
        return nil, 0, err
    }
    err = r.db.Select(&orders, query, userID, pageSize, (page-1)*pageSize)
    if err != nil {
        return nil, 0, err
    }
    for _, o := range orders {
        err = r.db.Select(&o.Items, "SELECT * FROM order_items WHERE order_id = $1", o.ID)
        if err != nil {
            return nil, 0, err
        }
    }
    return orders, total, nil
}