package polycode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

var serviceClient *ServiceClient = nil
var appConfig AppConfig = nil
var httpHandler http.Handler = nil

func init() {
	serviceClient = NewServiceClient("http://127.0.0.1:9999")
	appConfig = loadAppConfig()
}

type TaskStartEvent struct {
	Id         string    `json:"id"`
	SessionId  string    `json:"sessionId"`
	EntryPoint string    `json:"entryPoint"`
	Input      TaskInput `json:"input"`
}

type TaskCompleteEvent struct {
	Output         TaskOutput
	TaskInProgress bool
}

func loadAppConfig() AppConfig {
	// Load the YAML file
	yamlFile := "application.yml"
	var yamlData interface{}

	data, err := os.ReadFile(yamlFile)
	if os.IsNotExist(err) {
		log.Println("application.yml not found. Generating empty config...")
		yamlData = make(map[string]interface{}) // Create an empty config
	} else if err != nil {
		fmt.Printf("error reading yml file: %v\n", err)
		panic(err)
	} else {
		// Parse the YAML file into a map
		err = yaml.Unmarshal(data, &yamlData)
		if err != nil {
			fmt.Printf("error unmarshalling yml: %v\n", err)
			panic(err)
		}
	}

	yamlData = ConvertMap(yamlData)
	return yamlData.(map[string]interface{})
}

func Start(params ...any) {
	if len(params) == 1 {
		g, ok := params[1].(*gin.Engine)
		if ok {
			httpHandler = g.Handler()
		}
	}

	// Create a Gin router
	r := gin.Default()

	r.POST("/v1/invoke/api", invokeApiHandler)
	r.POST("/v1/invoke/service", invokeServiceHandler)

	c := make(chan bool)
	go func(c chan bool) {
		// Start the Gin server
		err := r.Run(fmt.Sprintf(":9998"))
		if err != nil {
			panic(err)
		}
		c <- true
	}(c)

	req := StartAppRequest{
		ClientPort: 9998,
	}

	println("starting app")
	err := serviceClient.StartApp(req)
	if err != nil {
		panic(err)
	}

	<-c
}

func invokeApiHandler(c *gin.Context) {
	input := c.Request
	output := runTask(nil, input)
	c.JSON(http.StatusOK, output)
}

func invokeServiceHandler(c *gin.Context) {
	var input TaskStartEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	output := runTask(nil, input)
	c.JSON(http.StatusOK, output)
}

func runTask(ctx context.Context, event any) (evt *TaskCompleteEvent) {
	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			// Check if it's the specific error
			if err, ok := r.(error); ok {
				if err == ErrTaskInProgress {
					fmt.Printf("task in progress\n")
					evt = &TaskCompleteEvent{TaskInProgress: true}
				} else {
					fmt.Printf("error %s\n", err.Error())
					evt = errorToTaskComplete(err)
				}
			} else {
				fmt.Printf("error %s\n", err.Error())
				evt = errorToTaskComplete(ErrUnknownError)
			}
		}
	}()

	switch it := event.(type) {
	case *TaskStartEvent:
		{
			service, err := GetService()
			if err != nil {
				return errorToTaskComplete(err)
			}

			inputObj, err := service.GetInputType(it.EntryPoint)
			if err != nil {
				return errorToTaskComplete(err)
			}
			err = json.Unmarshal([]byte(it.Input.TargetReq), inputObj)
			if err != nil {
				return errorToTaskComplete(err)
			}

			isWorkflow := service.IsWorkflow(it.EntryPoint)

			var ret any
			if isWorkflow {
				workflowCtx := WorkflowContext{
					ctx:           ctx,
					sessionId:     it.SessionId,
					serviceClient: serviceClient,
					config:        appConfig,
				}
				ret, err = service.ExecuteWorkflow(workflowCtx, it.EntryPoint, inputObj)
			} else {
				srvCtx := ServiceContext{
					ctx:       ctx,
					sessionId: it.SessionId,
					dataStore: NewDatabase(serviceClient, it.SessionId),
					config:    appConfig,
				}
				ret, err = service.ExecuteService(srvCtx, it.EntryPoint, inputObj)
			}

			if err != nil {
				return errorToTaskComplete(err)
			}

			output := TaskOutput{}
			if ret == nil {
				output.IsNull = true
			} else {
				output.Output = ret
			}
			return outputToTaskComplete(output)
		}
	case *http.Request:
		{
			if httpHandler == nil {
				return errorToTaskComplete(ErrBadRequest)
			}

			sessionId := it.Header.Get("x-polycode-task-session-id")
			if sessionId == "" {
				return errorToTaskComplete(ErrBadRequest)
			}

			path := it.Header.Get("x-polycode-task-api-path")
			if path == "" {
				return errorToTaskComplete(ErrBadRequest)
			}

			workflowCtx := WorkflowContext{
				ctx:           ctx,
				sessionId:     sessionId,
				serviceClient: serviceClient,
				config:        appConfig,
			}
			wkfCtx := context.WithValue(it.Context(), "polycode.context", workflowCtx)
			it = it.WithContext(wkfCtx)
			it.URL.Path = path
			resp := invokeHandler(httpHandler, it)
			return &TaskCompleteEvent{Output: TaskOutput{IsAsync: false, IsNull: false, Output: resp, Error: nil}}
		}
	}

	return errorToTaskComplete(ErrBadRequest)
}

func errorToTaskComplete(err error) *TaskCompleteEvent {
	ret := ErrTaskExecError.Wrap(err)
	println(fmt.Sprintf("task completed with error, %v", ret))
	output := TaskOutput{IsAsync: false, IsNull: false, Error: &ret}
	return outputToTaskComplete(output)
}

func outputToTaskComplete(output TaskOutput) *TaskCompleteEvent {
	return &TaskCompleteEvent{
		Output: output,
	}
}
