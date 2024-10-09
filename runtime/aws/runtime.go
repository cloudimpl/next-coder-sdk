package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/client"
	"github.com/CloudImpl-Inc/next-coder-sdk/client/db"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"reflect"
)

type Runtime struct {
	h        http.Handler
	dbClient *db.Client
}

func (l *Runtime) AppConfig() polycode.AppConfig {
	// Load the YAML file
	yamlFile := "application.yml" // The path to your YAML file
	var yamlData interface{}

	data, err := os.ReadFile(yamlFile)
	if os.IsNotExist(err) {
		log.Println("application.yml not found. Generating empty config...")
		yamlData = make(map[string]interface{}) // Create an empty config
	} else if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		os.Exit(1)
	} else {
		// Parse the YAML file into a map
		err = yaml.Unmarshal(data, &yamlData)
		if err != nil {
			fmt.Printf("Error unmarshalling YAML: %v\n", err)
			os.Exit(1)
		}
	}

	// Convert map[interface{}]interface{} to map[string]interface{}
	yamlData = polycode.ConvertMap(yamlData)
	return yamlData.(map[string]interface{})
}

func (l *Runtime) GetRuntime() polycode.Runtime {
	return l
}

func (l *Runtime) Name() string {
	return "aws-lambda"
}

func (l *Runtime) Start(params []any) error {
	println(fmt.Sprintf("params len %d\n", len(params)))

	l.dbClient = db.NewClient("http://localhost:6666")

	// Loop through params and print each item, including its type
	for i, param := range params {
		fmt.Printf("Param[%d]: %v (Type: %s)\n", i, param, reflect.TypeOf(param).String())
	}

	if len(params) == 2 {
		port, ok := params[0].(int)
		if !ok {
			return fmt.Errorf("port must be a int")
		}
		h, ok := params[1].(http.Handler)
		if !ok {
			fmt.Printf("casting failed")
			return fmt.Errorf("handler must be a http.Handler")
		}
		g, ok := params[1].(*gin.Engine)
		if ok {
			g.GET("/@routers", GetRoutesHandler(g))
		}
		return l.listenAndServe(fmt.Sprintf("127.0.0.1:%d", port), h)
	} else {
		err := l.listenAndServe("", nil)
		return err
	}
}

func (l *Runtime) InvokeWorkflow(workflowContext polycode.WorkflowContext, input polycode.TaskInput) (any, error) {
	if l.h == nil {
		return nil, ErrHttpHandlerNotSet
	}
	evt := events.APIGatewayProxyRequest{}

	if err := json.Unmarshal([]byte(input.TargetReq), &evt); err != nil {
		return nil, err
	}

	ctx := context.WithValue(workflowContext, "polycode.context", workflowContext)
	r, err := gateway.NewRequest(ctx, evt)
	if err != nil {
		return nil, err
	}
	w := gateway.NewResponse()
	l.h.ServeHTTP(w, r)

	resp := w.End()
	println(fmt.Sprintf("workflow response %v\n", resp))
	return resp, nil
}

// Invoke Handler implementation
func (l *Runtime) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	fmt.Printf("data received to lambda function %s\n", string(payload))
	evt := client.Event{}

	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, err
	}

	runtimeContext := client.NewRuntimeContext(ctx, db.NewDatabase(l.dbClient, evt.Context.SessionId))

	ret, err := client.RunTask(runtimeContext, l, evt)
	if err != nil {
		return nil, err
	}

	return json.Marshal(ret)
}

func (l *Runtime) listenAndServe(addr string, h http.Handler) error {
	fmt.Printf("listening on %s %v\n", addr, h)
	l.h = h
	if l.h == nil {
		fmt.Printf("default handler set\n")
		l.h = http.DefaultServeMux
	}

	lambda.StartHandler(l)

	return nil
}
