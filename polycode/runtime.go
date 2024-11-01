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
	"runtime/debug"
	"strings"
)

var serviceClient = NewServiceClient("http://127.0.0.1:9999")
var appConfig = loadAppConfig()
var serviceMap = make(map[string]Service)
var httpHandler *gin.Engine = nil

type Service interface {
	GetName() string
	GetInputType(method string) (any, error)
	ExecuteService(ctx ServiceContext, method string, input any) (any, error)
	ExecuteWorkflow(ctx WorkflowContext, method string, input any) (any, error)
	IsWorkflow(method string) bool
}

func RegisterService(service Service) {
	fmt.Println("client: register service ", service.GetName())

	if serviceMap[service.GetName()] != nil {
		panic(fmt.Sprintf("client: service %s already registered", service.GetName()))
	}

	serviceMap[service.GetName()] = service
}

func RegisterApi(engine *gin.Engine) {
	fmt.Println("client: register api")

	if httpHandler != nil {
		panic("client: api already registered")
	}

	httpHandler = engine
}

func StartApp() {
	if len(os.Args) > 1 {
		println("client: run cli command")
		err := runCliCommand(os.Args[1:])
		if err != nil {
			panic(err)
		}
	} else {
		go startApiServer()
		sendStartApp()
		fmt.Printf("client: app %s started on port %d\n", GetClientEnv().AppName, GetClientEnv().AppPort)
		select {}
	}
}

func getService(serviceName string) (Service, error) {
	service := serviceMap[serviceName]
	if service == nil {
		return nil, fmt.Errorf("client: service %s not registered", serviceName)
	}
	return service, nil
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

func sendStartApp() {
	var services []ServiceData
	for name := range serviceMap {
		services = append(services, ServiceData{
			Name: name,
			// ToDo: Add task info
		})
	}

	req := StartAppRequest{
		AppName:  GetClientEnv().AppName,
		AppPort:  GetClientEnv().AppPort,
		Services: services,
		Routes:   loadRoutes(),
	}

	err := serviceClient.StartApp(req)
	if err != nil {
		panic(err)
	}
}

func loadRoutes() []RouteData {
	var routes = make([]RouteData, 0)
	if httpHandler != nil {
		for _, route := range httpHandler.Routes() {
			fmt.Printf("client: route found %s %s\n", route.Method, route.Path)

			routes = append(routes, RouteData{
				Method: route.Method,
				Path:   route.Path,
			})
		}
	}
	return routes
}

func runTask(ctx context.Context, event any) (evt TaskCompleteEvent) {
	println(fmt.Sprintf("client: run task with event %v", reflect.TypeOf(event)))

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			// Check if it's the specific error
			if err, ok := r.(error); ok {
				if errors.Is(err, ErrTaskInProgress) {
					println("client: task in progress")
					evt = TaskCompleteEvent{}
				} else {
					fmt.Printf("client: task completed with error %s\n", err.Error())
					stackTrace := string(debug.Stack())
					println(stackTrace)
					evt = ErrorToTaskComplete(err)
				}
			} else {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				stackTrace := string(debug.Stack())
				println(stackTrace)
				evt = ErrorToTaskComplete(err)
			}
		}
	}()

	switch it := event.(type) {
	case TaskStartEvent:
		{
			println("client: handle task start event")

			service, err := getService(it.ServiceName)
			if err != nil {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				return ErrorToTaskComplete(err)
			}

			inputObj, err := service.GetInputType(it.EntryPoint)
			if err != nil {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				return ErrorToTaskComplete(err)
			}
			err = json.Unmarshal([]byte(it.Input.TargetReq), inputObj)
			if err != nil {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				return ErrorToTaskComplete(err)
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

				println(fmt.Sprintf("client: service %s exec workflow %s with session id %s", it.ServiceName,
					it.EntryPoint, it.SessionId))
				ret, err = service.ExecuteWorkflow(workflowCtx, it.EntryPoint, inputObj)
			} else {
				srvCtx := ServiceContext{
					ctx:       ctx,
					sessionId: it.SessionId,
					dataStore: NewDatabase(serviceClient, it.SessionId),
					config:    appConfig,
				}

				println(fmt.Sprintf("client: service %s exec handler %s with session id %s", it.ServiceName,
					it.EntryPoint, it.SessionId))
				ret, err = service.ExecuteService(srvCtx, it.EntryPoint, inputObj)
			}

			if err != nil {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				return ErrorToTaskComplete(err)
			}

			if ret == nil {
				println("client: task completed")
				return NilValueToTaskComplete()
			} else {
				println("client: task completed")
				return ValueToTaskComplete(ret)
			}
		}
	case ApiStartEvent:
		{
			println("client: handle http request")
			if httpHandler == nil {
				fmt.Printf("client: task completed with error %s\n", ErrBadRequest.Error())
				return ErrorToTaskComplete(ErrBadRequest)
			}

			apiCtx := ApiContext{
				ctx:           ctx,
				sessionId:     it.SessionId,
				serviceClient: serviceClient,
				config:        appConfig,
			}
			newCtx := context.WithValue(ctx, "polycode.context", apiCtx)

			req, err := ConvertToHttpRequest(newCtx, it.Request)
			if err != nil {
				fmt.Printf("client: task completed with error %s\n", err.Error())
				return ErrorToTaskComplete(err)
			}

			resp := invokeHandler(httpHandler, req)
			println("client: task completed")
			return ValueToTaskComplete(resp)
		}
	default:
		{
			fmt.Printf("client: invalid event type %v\n", reflect.TypeOf(event))
			fmt.Printf("client: task completed with error %s\n", ErrBadRequest.Error())
			return ErrorToTaskComplete(ErrBadRequest)
		}
	}
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

	println("client: create http request with workflow context")
	req, err := http.NewRequestWithContext(ctx, apiReq.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range apiReq.Header {
		req.Header.Set(key, value)
	}

	return req, nil
}

func ValueToTaskComplete(val any) TaskCompleteEvent {
	output := TaskOutput{
		Output: val,
	}
	return TaskCompleteEvent{
		Output: output,
	}
}

func NilValueToTaskComplete() TaskCompleteEvent {
	output := TaskOutput{
		IsNull: true,
	}
	return TaskCompleteEvent{
		Output: output,
	}
}

func ErrorToTaskComplete(err error) TaskCompleteEvent {
	ret := ErrTaskExecError.Wrap(err)
	output := TaskOutput{
		IsError: true,
		Error:   ret,
	}
	return TaskCompleteEvent{
		Output: output,
	}
}
