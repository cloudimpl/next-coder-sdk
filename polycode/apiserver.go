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
	println("api task started")
	var input ApiStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		errorOutput := ErrorEvent{
			Error: ErrInternal.Wrap(err),
		}
		fmt.Printf("api task failed %s", err.Error())
		c.JSON(http.StatusInternalServerError, errorOutput)
	} else {
		output := runApi(c, input)
		println("api task success")
		c.JSON(http.StatusOK, output)
	}
}

func invokeServiceHandler(c *gin.Context) {
	println("service task started")
	var input ServiceStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		errorOutput := ErrorEvent{
			Error: ErrInternal.Wrap(err),
		}
		fmt.Printf("service task failed %s", err.Error())
		c.JSON(http.StatusInternalServerError, errorOutput)
	} else {
		output := runService(c, input)
		println("service task success")
		c.JSON(http.StatusOK, output)
	}
}
