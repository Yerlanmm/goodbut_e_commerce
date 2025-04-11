package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func main() {
	r := gin.Default()

	// Inventory routes
	r.Any("/products/*proxyPath", proxyRequest("http://localhost:8001"))
	// Order routes
	r.Any("/orders/*proxyPath", proxyRequest("http://localhost:8002"))

	// Optional root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API Gateway!"})
	})

	r.Run(":8000") // The gateway runs on port 8000
}

// Generic proxy handler
func proxyRequest(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := target + c.Request.URL.Path
		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Forward the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach service"})
			return
		}
		defer resp.Body.Close()

		// Copy response
		c.Status(resp.StatusCode)
		for k, vv := range resp.Header {
			for _, v := range vv {
				c.Writer.Header().Add(k, v)
			}
		}
		io.Copy(c.Writer, resp.Body)
	}
}
