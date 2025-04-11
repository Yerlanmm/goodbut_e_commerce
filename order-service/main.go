package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"html/template"
	"strconv"
)


type Order struct {
	ID          int      `json:"id"`
	ProductIDs  []int    `json:"product_ids"`
	TotalAmount float64  `json:"total_amount"`
	Status      string   `json:"status"`
}

var orders = []Order{
	{ID: 1, ProductIDs: []int{1, 2}, TotalAmount: 30.5, Status: "Pending"},
	{ID: 2, ProductIDs: []int{2}, TotalAmount: 20.0, Status: "Completed"},
}

func main() {
	r := gin.Default()

	
	r.POST("/orders", func(c *gin.Context) {
		var newOrder Order
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newOrder.ID = len(orders) + 1
		orders = append(orders, newOrder)
		c.JSON(http.StatusCreated, newOrder)
	})

	
	r.GET("/orders/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for _, o := range orders {
			if o.ID == id {
				c.JSON(http.StatusOK, o)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	
	r.PATCH("/orders/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var updatedOrder Order
		if err := c.ShouldBindJSON(&updatedOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, o := range orders {
			if o.ID == id {
				orders[i].Status = updatedOrder.Status
				c.JSON(http.StatusOK, orders[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	
	r.DELETE("/orders/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for i, o := range orders {
			if o.ID == id {
				orders = append(orders[:i], orders[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	})

	
	r.GET("/orders", func(c *gin.Context) {
		tmpl, err := template.New("orders").Parse(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Orders</title>
				<style>
					.card {
						border: 1px solid #ddd;
						padding: 10px;
						margin: 10px;
						border-radius: 5px;
					}
					.card h3 {
						margin: 0;
					}
					.card button {
						margin-top: 10px;
					}
				</style>
			</head>
			<body>
				<h1>Order List</h1>
				{{range .}}
					<div class="card">
						<h3>Order ID: {{.ID}}</h3>
						<p>Status: {{.Status}}</p>
						<p>Products: {{.ProductIDs}}</p>
						<p>Total Amount: ${{.TotalAmount}}</p>
						<form action="/orders/{{.ID}}" method="POST">
							<button type="submit" name="_method" value="DELETE">Delete</button>
						</form>
					</div>
				{{else}}
					<p>No orders available.</p>
				{{end}}
			</body>
			</html>
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render HTML"})
			return
		}
		c.Header("Content-Type", "text/html")
		tmpl.Execute(c.Writer, orders)
	})

	
	r.Run(":8002")
}
