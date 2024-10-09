package client

import (
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"os"
)

var serviceClient *ServiceClient = nil

func init() {
}

func getServiceClient() *ServiceClient {
	if serviceClient == nil {
		serviceClient = NewServiceClient("http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	}
	return serviceClient
}

func RunTask(ctx RuntimeContext, runtime RuntimeSupport, event Event) (*TaskCompleteEvent, error) {
	return _RunTask(ctx, runtime, event)
}

func _RunTask(ctx RuntimeContext, runtime RuntimeSupport, event Event) (evt *TaskCompleteEvent, errRet error) {

	defer func() {
		// Recover from panic and check for a specific error
		if r := recover(); r != nil {
			// Check if it's the specific error
			if err, ok := r.(error); ok {
				if polycode.IsError(err, ErrTaskInProgress) {
					fmt.Printf("task in progress\n")
					evt = &TaskCompleteEvent{TaskInProgress: true}
					errRet = nil
				} else {
					println("panic not recovered")
					panic(r)
				}
			} else {
				// If it's not the specific error, re-panic
				panic(r)
			}
		}
	}()

	if event.Error != "" {
		return ConvertToCompletionEvent(fmt.Errorf(event.Error)), nil
	}
	if event.Type == InternalCall {
		output := polycode.TaskOutput{
			IsAsync: false,
			IsNull:  false,
			Output:  nil,
			Error:   nil,
		}
		if event.TaskInput.NoArg {
			output.IsNull = true
		} else {
			output.Output = event.TaskInput.TargetReq
		}

		return CreateTaskCompleteEvent(output, nil), nil
	}
	fmt.Printf("run task started %v\n", event)
	if event.Type == EventAPICall {
		workflowCtx := WorkflowContext{
			ctx:           ctx,
			sessionId:     event.Context.SessionId,
			serviceClient: getServiceClient(),
			config:        runtime.GetRuntime().AppConfig(),
		}
		ret, err := runtime.InvokeWorkflow(workflowCtx, event.TaskInput)
		if err != nil {
			return nil, err
		}
		output := polycode.TaskOutput{}
		if ret == nil {
			output.IsNull = true
		} else {
			output.Output = ret
		}
		return CreateTaskCompleteEvent(output, nil), nil
	} else if event.Type == EventServiceCall {
		db := ctx.GetDB()

		service, err := polycode.GetService()
		if err != nil {
			return ConvertToCompletionEvent(err), nil
		}

		inputObj, err := service.GetInputType(event.Context.EntryPoint)
		if err != nil {
			return ConvertToCompletionEvent(err), nil
		}
		err = json.Unmarshal([]byte(event.TaskInput.TargetReq), inputObj)
		if err != nil {
			return ConvertToCompletionEvent(err), nil
		}

		isWorkflow := service.IsWorkflow(event.Context.EntryPoint)

		var ret any
		if isWorkflow {
			workflowCtx := WorkflowContext{
				ctx:           ctx,
				sessionId:     event.Context.SessionId,
				serviceClient: getServiceClient(),
				config:        runtime.GetRuntime().AppConfig(),
			}
			ret, err = service.ExecuteWorkflow(workflowCtx, event.Context.EntryPoint, inputObj)
		} else {
			srvCtx := ServiceContext{
				ctx:       ctx,
				sessionId: event.Context.SessionId,
				db:        db,
				config:    runtime.GetRuntime().AppConfig(),
			}
			ret, err = service.ExecuteService(srvCtx, event.Context.EntryPoint, inputObj)
		}

		if err != nil {
			return ConvertToCompletionEvent(err), nil
		}

		output := polycode.TaskOutput{}
		if ret == nil {
			output.IsNull = true
		} else {
			output.Output = ret
		}
		return CreateTaskCompleteEvent(output, db.GetTx()), nil
	} else {
		return ConvertToCompletionEvent(ErrBadRequest), nil
	}

}

func ConvertToCompletionEvent(err error) *TaskCompleteEvent {
	ret := ErrTaskExecError.Wrap(err)
	println(fmt.Sprintf("task completed with error, %v", ret))
	return &TaskCompleteEvent{Output: polycode.TaskOutput{IsAsync: false, IsNull: false, Error: &ret}}
}
