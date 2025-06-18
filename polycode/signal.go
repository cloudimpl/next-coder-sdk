package polycode

import "fmt"

type Signal struct {
	name          string
	sessionId     string
	serviceClient *ServiceClient
}

func (s *Signal) Await() Response {
	req := SignalWaitRequest{
		SignalName: s.name,
	}

	output, err := s.serviceClient.WaitForSignal(s.sessionId, req)
	if err != nil {
		fmt.Printf("client: signal await error: %v\n", err)
		return Response{
			output:  nil,
			isError: true,
			error:   ErrTaskExecError.Wrap(err),
		}
	}

	fmt.Printf("client: signal await output: %v\n", output)
	return Response{
		output:  output.Output,
		isError: output.IsError,
		error:   output.Error,
	}
}

func (s *Signal) EmitValue(taskId string, data any) error {
	req := SignalEmitRequest{
		TaskId:     taskId,
		SignalName: s.name,
		Output:     data,
		IsError:    false,
	}

	return s.serviceClient.EmitSignal(s.sessionId, req)
}

func (s *Signal) EmitError(taskId string, err Error) error {
	req := SignalEmitRequest{
		TaskId:     taskId,
		SignalName: s.name,
		IsError:    true,
		Error:      err,
	}

	return s.serviceClient.EmitSignal(s.sessionId, req)
}
