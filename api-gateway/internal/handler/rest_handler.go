package handler

import (
	pbInventory "E-Commerce/inventory-service/proto"
	pbOrder "E-Commerce/order-service/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RESTHandler struct {
	inventoryClient pbInventory.InventoryServiceClient
	orderClient     pbOrder.OrderServiceClient
}

func NewRESTHandler(inventoryClient pbInventory.InventoryServiceClient, orderClient pbOrder.OrderServiceClient) *RESTHandler {
	return &RESTHandler{
		inventoryClient: inventoryClient,
		orderClient:     orderClient,
	}
}

func (h *RESTHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		Stock       int     `json:"stock" binding:"required"`
		CategoryID  string  `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.inventoryClient.CreateProduct(c.Request.Context(), &pbInventory.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       float32(req.Price),
		Stock:       int32(req.Stock),
		CategoryId:  req.CategoryID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp.Product)
}

func (h *RESTHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.inventoryClient.GetProduct(c.Request.Context(), &pbInventory.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Product)
}

func (h *RESTHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		CategoryID  string  `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.inventoryClient.UpdateProduct(c.Request.Context(), &pbInventory.UpdateProductRequest{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       float32(req.Price),
		Stock:       int32(req.Stock),
		CategoryId:  req.CategoryID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Product)
}

func (h *RESTHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.inventoryClient.DeleteProduct(c.Request.Context(), &pbInventory.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

func (h *RESTHandler) ListProducts(c *gin.Context) {
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	resp, err := h.inventoryClient.ListProducts(c.Request.Context(), &pbInventory.ListProductsRequest{
		CategoryId: categoryID,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
		"total":    resp.Total,
	})
}

func (h *RESTHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Items []struct {
			ProductID string `json:"product_id" binding:"required"`
			Quantity  int    `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pbItems := make([]*pbOrder.OrderItem, len(req.Items))
	for i, item := range req.Items {
		pbItems[i] = &pbOrder.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		}
	}

	resp, err := h.orderClient.CreateOrder(c.Request.Context(), &pbOrder.CreateOrderRequest{
		UserId: userID.(string),
		Items:  pbItems,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp.Order)
}

func (h *RESTHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.orderClient.GetOrder(c.Request.Context(), &pbOrder.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Order)
}

func (h *RESTHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.orderClient.UpdateOrder(c.Request.Context(), &pbOrder.UpdateOrderRequest{
		Id:     id,
		Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp.Order)
}

func (h *RESTHandler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	queryUserID := c.Query("user_id")

	if role == "admin" && queryUserID != "" {
		userID = queryUserID
	} else {
		userID = userID.(string)
	}

	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	resp, err := h.orderClient.ListOrders(c.Request.Context(), &pbOrder.ListOrdersRequest{
		UserId:   userID.(string),
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"orders": resp.Orders,
		"total":  resp.Total,
	})
}
