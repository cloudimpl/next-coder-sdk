package api

import (
	"github.com/cloudimpl/next-coder-sdk/apicontext"
	"github.com/cloudimpl/next-coder-sdk/polycode"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FromWorkflow[Input any, Output any](f func(polycode.WorkflowContext, Input) (Output, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		apiCtx, err := apicontext.FromContext(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to execute workflow: " + err.Error(),
			})
			return
		}

		var input Input
		if err = c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}

		err = polycode.GetValidator().Validate(input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request",
			})
			return
		}

		output, err := f(apiCtx, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to execute workflow: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
