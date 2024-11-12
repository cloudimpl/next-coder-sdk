package polycode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	log.Println("client: api request received")
	var input ApiStartEvent

	if err := c.ShouldBindJSON(&input); err != nil {
		errorOutput := ErrorEvent{
			Error: ErrInternal.Wrap(err),
		}

		c.JSON(http.StatusInternalServerError, errorOutput)
		log.Println("client: api request error")
	} else {
		output := runApi(c, input)
		c.JSON(http.StatusOK, output)
		log.Println("client: api request completed")
	}
}

func invokeServiceHandler(c *gin.Context) {
	log.Println("client: service request received")
	var input ServiceStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		errorOutput := ErrorEvent{
			Error: ErrInternal.Wrap(err),
		}

		c.JSON(http.StatusInternalServerError, errorOutput)
		log.Println("client: service request error")
	} else {
		output := runService(c, input)
		c.JSON(http.StatusOK, output)
		log.Println("client: service request completed")
	}
}
