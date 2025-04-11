package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"html/template"
	"strconv"
)


type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

var products = []Product{
	{ID: 1, Name: "Product 1", Description: "Description 1", Price: 10.5, Stock: 100},
	{ID: 2, Name: "Product 2", Description: "Description 2", Price: 20.0, Stock: 50},
}

func main() {
	r := gin.Default()

	
	r.POST("/products", func(c *gin.Context) {
		var newProduct Product
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newProduct.ID = len(products) + 1
		products = append(products, newProduct)
		c.JSON(http.StatusCreated, newProduct)
	})

	
	r.GET("/products/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for _, p := range products {
			if p.ID == id {
				c.JSON(http.StatusOK, p)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	
	r.PATCH("/products/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var updatedProduct Product
		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, p := range products {
			if p.ID == id {
				products[i] = updatedProduct
				products[i].ID = id 
				c.JSON(http.StatusOK, products[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})

	
	r.DELETE("/products/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for i, p := range products {
			if p.ID == id {
				products = append(products[:i], products[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	})


	r.GET("/products", func(c *gin.Context) {
		tmpl, err := template.New("index").Parse(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Products</title>
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
				<h1>Product List</h1>
				{{range .}}
					<div class="card">
						<h3>{{.Name}}</h3>
						<p>{{.Description}}</p>
						<p><strong>Price:</strong> ${{.Price}}</p>
						<p><strong>Stock:</strong> {{.Stock}}</p>
						<form action="/products/{{.ID}}" method="POST">
							<button type="submit" name="_method" value="DELETE">Delete</button>
						</form>
					</div>
				{{else}}
					<p>No products available.</p>
				{{end}}
			</body>
			</html>
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render HTML"})
			return
		}
		c.Header("Content-Type", "text/html")
		tmpl.Execute(c.Writer, products)
	})

	r.Run(":8001") 
}
