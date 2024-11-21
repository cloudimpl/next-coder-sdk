package polycode

import (
	"context"
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
		fmt.Printf("client: service %s already registered\n", service.GetName())
	} else {
		serviceMap[service.GetName()] = service
	}
}

func StartApp(args ...any) {
	if len(args) > 1 {
		log.Fatal("client: invalid start app arguments")
	}

	if len(args) > 0 {
		g, ok := args[0].(*gin.Engine)
		if !ok {
			log.Fatalf("client: invalid start app arguments")
		}

		log.Println("client: register api")
		httpHandler = g
	}

	if len(os.Args) > 1 {
		log.Println("client: run cli command")
		err := runCliCommand(os.Args[1:])
		if err != nil {
			log.Fatalf("client: %s\n", err.Error())
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

	var err error
	for {
		err = serviceClient.StartApp(req)
		if err == nil {
			break
		}
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

func runService(ctx context.Context, taskLogger Logger, event ServiceStartEvent) (evt ServiceCompleteEvent) {
	taskLogger.Info().Msg(fmt.Sprintf("service started %s.%s", event.Service, event.Method))

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			recovered, ok := r.(error)

			if ok && errors.Is(recovered, ErrTaskStopped) {
				taskLogger.Info().Msg("service stopped")
				evt = ValueToServiceComplete(nil)
			} else {
				stackTrace := string(debug.Stack())
				fmt.Printf("stack trace %s\n", stackTrace)

				err2 := ErrInternal.Wrap(fmt.Errorf("recovered type %T", r))
				taskLogger.Error().Msg(err2.Error())
				evt = ErrorToServiceComplete(err2)
			}
		}
	}()

	service, err := getService(event.Service)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2)
	}

	inputObj, err := service.GetInputType(event.Method)
	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2)
	}

	err = ConvertType(event.Input, inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2)
	}

	err = currentValidator.Validate(inputObj)
	if err != nil {
		err2 := ErrBadRequest.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2)
	}

	ctxImpl := &ContextImpl{
		ctx:           ctx,
		sessionId:     event.SessionId,
		dataStore:     NewDatabase(serviceClient, event.SessionId),
		fileStore:     NewFileStore(serviceClient, event.SessionId),
		config:        appConfig,
		serviceClient: serviceClient,
		logger:        taskLogger,
	}

	var ret any
	if service.IsWorkflow(event.Method) {
		taskLogger.Info().Msg(fmt.Sprintf("service %s exec workflow %s with session id %s", event.Service,
			event.Method, event.SessionId))
		ret, err = service.ExecuteWorkflow(ctxImpl, event.Method, inputObj)
	} else {
		taskLogger.Info().Msg(fmt.Sprintf("service %s exec handler %s with session id %s", event.Service,
			event.Method, event.SessionId))
		ret, err = service.ExecuteService(ctxImpl, event.Method, inputObj)
	}

	if err != nil {
		err2 := ErrServiceExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToServiceComplete(err2)
	}

	taskLogger.Info().Msg("service completed")
	return ValueToServiceComplete(ret)
}

func runApi(ctx context.Context, taskLogger Logger, event ApiStartEvent) (evt ApiCompleteEvent) {
	taskLogger.Info().Msg(fmt.Sprintf("api started %s %s", event.Request.Method, event.Request.Path))

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			recovered, ok := r.(error)

			if ok && errors.Is(recovered, ErrTaskStopped) {
				taskLogger.Info().Msg("api stopped")
				evt = ApiCompleteEvent{
					Response: ApiResponse{
						StatusCode:      202,
						Header:          make(map[string]string),
						Body:            "",
						IsBase64Encoded: false,
					},
				}
			} else {
				err2 := ErrInternal.Wrap(recovered)
				stackTrace := string(debug.Stack())
				taskLogger.Error().Msg(err2.Error())
				taskLogger.Error().Msg(fmt.Sprintf("stack trace %s", stackTrace))
				evt = ApiCompleteEvent{
					Response: ApiResponse{
						StatusCode:      500,
						Header:          make(map[string]string),
						Body:            err2.ToJson(),
						IsBase64Encoded: false,
					},
				}
			}
		}
	}()

	if httpHandler == nil {
		err2 := ErrApiExecError.Wrap(errors.New("http handler not set"))
		taskLogger.Error().Msg(err2.Error())
		return ErrorToApiComplete(err2)
	}

	ctxImpl := &ContextImpl{
		ctx:           ctx,
		sessionId:     event.SessionId,
		dataStore:     NewDatabase(serviceClient, event.SessionId),
		fileStore:     NewFileStore(serviceClient, event.SessionId),
		config:        appConfig,
		serviceClient: serviceClient,
		logger:        taskLogger,
	}

	newCtx := context.WithValue(ctx, "polycode.context", ctxImpl)
	httpReq, err := ConvertToHttpRequest(newCtx, event.Request)
	if err != nil {
		err2 := ErrApiExecError.Wrap(err)
		taskLogger.Error().Msg(err2.Error())
		return ErrorToApiComplete(err2)
	}

	res := ManualInvokeHandler(httpHandler, httpReq)
	taskLogger.Info().Msg("api completed")
	return ApiCompleteEvent{
		Response: res,
	}
}

func ValueToServiceComplete(output any) ServiceCompleteEvent {
	return ServiceCompleteEvent{
		Output:  output,
		IsError: false,
		Error:   Error{},
	}
}

func ErrorToServiceComplete(err Error) ServiceCompleteEvent {
	return ServiceCompleteEvent{
		Output:  nil,
		IsError: true,
		Error:   err,
	}
}

func ErrorToApiComplete(err Error) ApiCompleteEvent {
	return ApiCompleteEvent{
		Response: ApiResponse{
			StatusCode:      500,
			Header:          make(map[string]string),
			Body:            err.ToJson(),
			IsBase64Encoded: false,
		},
	}
}
