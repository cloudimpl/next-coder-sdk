package polycode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"runtime/debug"
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
	log.Println("client: register service ", service.GetName())

	if serviceMap[service.GetName()] != nil {
		panic(fmt.Sprintf("client: service %s already registered", service.GetName()))
	}

	serviceMap[service.GetName()] = service
}

func RegisterApi(engine *gin.Engine) {
	log.Println("client: register api")

	if httpHandler != nil {
		panic("client: api already registered")
	}

	httpHandler = engine
}

func StartApp() {
	if len(os.Args) > 1 {
		log.Println("client: run cli command")
		err := runCliCommand(os.Args[1:])
		if err != nil {
			panic(err)
		}
	} else {
		go startApiServer()
		sendStartApp()
		log.Printf("client: app %s started on port %d\n", GetClientEnv().AppName, GetClientEnv().AppPort)
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
		log.Println("client: application.yml not found. generating empty config")
		yamlData = make(map[string]interface{}) // Create an empty config
	} else if err != nil {
		log.Printf("client: error reading yml file: %v\n", err)
		panic(err)
	} else {
		// Parse the YAML file into a map
		err = yaml.Unmarshal(data, &yamlData)
		if err != nil {
			log.Printf("client: error unmarshalling yml: %v\n", err)
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
			log.Printf("client: route found %s %s\n", route.Method, route.Path)

			routes = append(routes, RouteData{
				Method: route.Method,
				Path:   route.Path,
			})
		}
	}
	return routes
}

func runTask(ctx context.Context, event TaskStartEvent) (evt TaskCompleteEvent) {
	log.Println("client: handle task start event")
	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			err, ok := r.(error)

			if ok && errors.Is(err, ErrTaskInProgress) {
				log.Println("client: task in progress")
				evt = ValueToTaskComplete(nil)
			}

			log.Printf("client: task completed with error %s\n", err.Error())
			stackTrace := string(debug.Stack())
			println(stackTrace)
			evt = ErrorToTaskComplete(err)
		}
	}()

	service, err := getService(event.ServiceName)
	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	inputObj, err := service.GetInputType(event.EntryPoint)
	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	err = ConvertType(event.Input, inputObj)
	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	isWorkflow := service.IsWorkflow(event.EntryPoint)

	var ret any
	if isWorkflow {
		workflowCtx := WorkflowContext{
			ctx:           ctx,
			sessionId:     event.SessionId,
			serviceClient: serviceClient,
			config:        appConfig,
		}

		println(fmt.Sprintf("client: service %s exec workflow %s with session id %s", event.ServiceName,
			event.EntryPoint, event.SessionId))
		ret, err = service.ExecuteWorkflow(workflowCtx, event.EntryPoint, inputObj)
	} else {
		srvCtx := ServiceContext{
			ctx:       ctx,
			sessionId: event.SessionId,
			dataStore: NewDatabase(serviceClient, event.SessionId),
			config:    appConfig,
		}

		println(fmt.Sprintf("client: service %s exec handler %s with session id %s", event.ServiceName,
			event.EntryPoint, event.SessionId))
		ret, err = service.ExecuteService(srvCtx, event.EntryPoint, inputObj)
	}

	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	if ret == nil {
		println("client: task completed")
		return ValueToTaskComplete(nil)
	} else {
		retJson, err := json.Marshal(ret)
		if err != nil {
			fmt.Printf("client: task completed with error %s\n", err.Error())
			return ErrorToTaskComplete(err)
		}

		println("client: task completed")
		return ValueToTaskComplete(string(retJson))
	}
}

func runApi(ctx context.Context, event ApiStartEvent) (evt TaskCompleteEvent) {
	log.Println("client: handle http request")
	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			err, ok := r.(error)

			if ok && errors.Is(err, ErrTaskInProgress) {
				log.Println("client: api in progress")
				evt = ValueToTaskComplete(nil)
			}

			log.Printf("client: api completed with error %s\n", err.Error())
			stackTrace := string(debug.Stack())
			println(stackTrace)
			evt = ErrorToTaskComplete(err)
		}
	}()

	if httpHandler == nil {
		log.Printf("client: task completed with error %s\n", ErrBadRequest.Error())
		return ErrorToTaskComplete(ErrBadRequest)
	}

	apiCtx := ApiContext{
		ctx:           ctx,
		sessionId:     event.SessionId,
		serviceClient: serviceClient,
		config:        appConfig,
	}
	newCtx := context.WithValue(ctx, "polycode.context", apiCtx)

	res, err := invokeHandler(newCtx, httpHandler, event.Request)
	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("client: task completed with error %s\n", err.Error())
		return ErrorToTaskComplete(err)
	}

	println("client: task completed")
	return ValueToTaskComplete(string(resJson))
}

func ValueToTaskComplete(output any) TaskCompleteEvent {
	taskOutput := TaskOutput{
		Output:  output,
		IsError: false,
		Error:   Error{},
	}

	return TaskCompleteEvent{
		Output: taskOutput,
	}
}

func ErrorToTaskComplete(err error) TaskCompleteEvent {
	taskOutput := TaskOutput{
		Output:  nil,
		IsError: true,
		Error:   ErrTaskExecError.Wrap(err),
	}

	return TaskCompleteEvent{
		Output: taskOutput,
	}
}
