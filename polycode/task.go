package polycode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var serviceClient *ServiceClient = nil
var appConfig AppConfig = nil
var httpHandler http.Handler = nil

func init() {
	serviceClient = NewServiceClient("http://127.0.0.1:9999")
	appConfig = loadAppConfig()
}

type RouteData struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type TaskStartEvent struct {
	Id         string    `json:"id"`
	SessionId  string    `json:"sessionId"`
	EntryPoint string    `json:"entryPoint"`
	Input      TaskInput `json:"input"`
}

type TaskCompleteEvent struct {
	Output TaskOutput
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

func startApiServer(port int) {
	// Create a Gin router
	r := gin.Default()

	r.GET("/v1/health", invokeHealthCheck)
	r.POST("/v1/invoke/api", invokeApiHandler)
	r.POST("/v1/invoke/service", invokeServiceHandler)

	// Start the Gin server
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func sendStartApp(port int, routes []RouteData) {
	req := StartAppRequest{
		ClientPort: port,
		Routes:     routes,
	}

	err := serviceClient.StartApp(req)
	if err != nil {
		panic(err)
	}
}

func LoadRoutes(engine *gin.Engine) []RouteData {
	var routes []RouteData
	for _, route := range engine.Routes() {
		routes = append(routes, RouteData{
			Method: route.Method,
			Path:   route.Path,
		})
	}
	return routes
}

func Start(params ...any) {
	var routes []RouteData
	if len(params) == 1 {
		g, ok := params[0].(*gin.Engine)
		if ok {
			httpHandler = g.Handler()
			LoadRoutes(g)
		}
	}

	go startApiServer(9998)
	sendStartApp(9998, routes)
	println("client: app started")

	select {}
}

func ConvertToHttpRequest(ctx context.Context, apiReq ApiRequest) (*http.Request, error) {
	// Build the URL
	url := apiReq.Path
	if len(apiReq.Query) > 0 {
		queryParams := "?"
		for key, value := range apiReq.Query {
			queryParams += key + "=" + value + "&"
		}
		queryParams = strings.TrimSuffix(queryParams, "&")
		url += queryParams
	}

	// Create a new HTTP request
	var body io.Reader
	if apiReq.Body != "" {
		body = bytes.NewReader([]byte(apiReq.Body))
	} else {
		body = nil
	}

	req, err := http.NewRequest(apiReq.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range apiReq.Header {
		req.Header.Set(key, value)
	}

	req.WithContext(ctx)
	return req, nil
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

func runTask(ctx context.Context, event any) (evt *TaskCompleteEvent) {
	println(fmt.Sprintf("client: run task with event %v", reflect.TypeOf(event)))

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			// Check if it's the specific error
			if err, ok := r.(error); ok {
				if errors.Is(err, ErrTaskInProgress) {
					fmt.Printf("task in progress\n")
					evt = &TaskCompleteEvent{}
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
			println("client: handle task start event")

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

				println(fmt.Sprintf("client: exec workflow %s with session id %s", it.EntryPoint, it.SessionId))
				ret, err = service.ExecuteWorkflow(workflowCtx, it.EntryPoint, inputObj)
			} else {
				srvCtx := ServiceContext{
					ctx:       ctx,
					sessionId: it.SessionId,
					dataStore: NewDatabase(serviceClient, it.SessionId),
					config:    appConfig,
				}

				println(fmt.Sprintf("client: exec service %s with session id %s", it.EntryPoint, it.SessionId))
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

			println("client: task completed")
			return outputToTaskComplete(output)
		}
	case *ApiStartEvent:
		{
			println("client: handle http request")

			if httpHandler == nil {
				println("client: http handler not found")
				return errorToTaskComplete(ErrBadRequest)
			}

			workflowCtx := WorkflowContext{
				ctx:           ctx,
				sessionId:     it.SessionId,
				serviceClient: serviceClient,
				config:        appConfig,
			}
			wkfCtx := context.WithValue(ctx, "polycode.context", workflowCtx)

			req, err := ConvertToHttpRequest(wkfCtx, it.Request)
			if err != nil {
				println("client: failed to convert api request")
				return errorToTaskComplete(err)
			}

			resp := invokeHandler(httpHandler, req)
			println("client: task completed")
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
