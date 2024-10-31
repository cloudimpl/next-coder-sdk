package polycode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func startApiServer() {
	// Create a Gin router
	r := gin.Default()

	r.GET("/v1/health", invokeHealthCheck)
	r.POST("/v1/invoke/api", invokeApiHandler)
	r.POST("/v1/invoke/service", invokeServiceHandler)

	// Start the Gin server
	err := r.Run(fmt.Sprintf(":%d", GetClientEnv().AppPort))
	if err != nil {
		panic(err)
	}
}

func invokeHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func invokeApiHandler(c *gin.Context) {
	println("client: api request received")
	var input ApiStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	output := runTask(c, &input)
	println("client: api request completed")
	c.JSON(http.StatusOK, output)
}

func invokeServiceHandler(c *gin.Context) {
	println("client: service request received")
	var input TaskStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	output := runTask(c, &input)
	println("client: service request completed")
	c.JSON(http.StatusOK, output)
}
