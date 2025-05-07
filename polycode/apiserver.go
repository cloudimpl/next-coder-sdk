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
		log.Fatalf("Failed to start api server: %s", err.Error())
	}
}

func invokeHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func invokeApiHandler(c *gin.Context) {
	taskLogger := CreateLogger("task")

	var input ApiStartEvent
	var output ApiCompleteEvent

	taskLogger.Info().Msg("api task started")
	if err := c.ShouldBindJSON(&input); err != nil {
		output = ErrorToApiComplete(ErrInternal.Wrap(err))
		taskLogger.Error().Msg(fmt.Sprintf("api task failed %s", err.Error()))
	} else {
		output = runApi(c, taskLogger, input)
		taskLogger.Info().Msg("api task success")
	}

	c.JSON(http.StatusOK, output)
}

func invokeServiceHandler(c *gin.Context) {
	taskLogger := CreateLogger("task")

	var input ServiceStartEvent
	var output ServiceCompleteEvent

	taskLogger.Info().Msg("service task started")
	if err := c.ShouldBindJSON(&input); err != nil {
		output = ErrorToServiceComplete(ErrInternal.Wrap(err))
		taskLogger.Error().Msg(fmt.Sprintf("service task failed %s", err.Error()))
	} else {
		output = runService(c, taskLogger, input)
		taskLogger.Info().Msg("service task success")
	}

	c.JSON(http.StatusOK, output)
}
