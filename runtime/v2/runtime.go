package v2

import (
	"github.com/CloudImpl-Inc/next-coder-sdk/client"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"github.com/gin-gonic/gin"
	"os"
)

type Runtime struct {
}

func (r Runtime) Name() string {
	return "V2"
}

func (r Runtime) AppConfig() polycode.AppConfig {
	return polycode.AppConfig{}
}

func (r *Runtime) OnRequest(c *gin.Context) {

	c.JSON(200, gin.H{
		"response": string("hello from client"),
	})
}

func (r Runtime) Start(params []any) error {
	serviceClient := client.NewServiceClient("http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	err := serviceClient.StartApp(client.StartAppRequest{})
	if err != nil {
		return err
	}
	//runtime = r
	//start gin server
	gin := gin.Default()
	gin.POST("/invoke", r.OnRequest)
	gin.Run(":8080")
	return nil
}
