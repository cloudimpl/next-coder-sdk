package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CloudImpl-Inc/next-coder-sdk/client/db"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
	"net/http"
	"time"
)

const (
	TaskPending TaskStatus = iota
	TaskRunning
	TaskSuccess
	TaskFailed
	TaskCancelled
)

type TaskStatus int8

type Request struct {
	Context RemoteTaskContext
	Input   polycode.TaskInput
}

type Output struct {
	SessionId string `json:"sessionId"`
}

type UrlCache struct {
	serviceIdUrl string
	//callTaskUrl  string
}

// ServiceClient is a reusable client for calling the service API
type ServiceClient struct {
	httpClient *http.Client
	baseURL    string
	urlCache   UrlCache
}

func CreateTaskCompleteEvent(output polycode.TaskOutput, tx *db.Tx) *TaskCompleteEvent {
	return &TaskCompleteEvent{
		Tx:     tx,
		Output: output,
	}
}

// NewServiceClient creates a new ServiceClient with a reusable HTTP client
func NewServiceClient(baseURL string) *ServiceClient {
	return &ServiceClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30, // Set a reasonable timeout for HTTP requests
		},
		baseURL: baseURL,
		urlCache: UrlCache{
			serviceIdUrl: fmt.Sprintf("%s/system/service/id", baseURL),
			//callTaskUrl:  fmt.Sprintf("%s/system/task/%s/%s/call", baseURL),
		},
	}
}

//func (sc *ServiceClient) startTask(context *TaskContext, input *polycode.TaskInput) (string, error) {
//
//	req := Request{
//		Context: context,
//		Input:   input,
//	}
//	b, err := json.Marshal(req)
//	if err != nil {
//		return "", fmt.Errorf("failed to marshal JSON: %w", err)
//	}
//	resp, err := sc.httpClient.Post(fmt.Sprintf("%s/system/task/init", sc.baseURL), "application/json", bytes.NewBuffer(b))
//	if err != nil {
//		return "", fmt.Errorf("failed to make HTTP request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Check if the response status code is 200 OK
//	if resp.StatusCode != http.StatusOK {
//		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	output := &Output{}
//	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
//		return "", fmt.Errorf("failed to decode JSON response: %w", err)
//	}
//	return output.SessionId, nil
//}

func (sc *ServiceClient) completeTask(sessionId string, completeEvent *TaskCompleteEvent) error {
	b, err := json.Marshal(completeEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	resp, err := sc.httpClient.Post(fmt.Sprintf("%s/system/tasks/%s/complete", sc.baseURL, sessionId), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (sc *ServiceClient) ExecTask(ctx RemoteTaskContext, method string, input polycode.TaskInput) (*polycode.TaskOutput, error) {

	req := Request{
		Context: ctx,
		Input:   input,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	// Make the HTTP POST request
	resp, err := sc.httpClient.Post(fmt.Sprintf("%s/v1/system/tasks/%s/exec", sc.baseURL, ctx.SessionId), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	output := &polycode.TaskOutput{}
	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if output.IsAsync {
		panic(ErrTaskInProgress) //task in progress panic
	}
	// Return the taskId from the response
	return output, nil
}

//func (sc *ServiceClient) dbCommit(sessionId string, request *TxRequest) error {
//
//	b, err := json.Marshal(request)
//	if err != nil {
//		return fmt.Errorf("failed to marshal JSON: %w", err)
//	}
//	resp, err := sc.httpClient.Post(fmt.Sprintf("%s/system/%s/db/commit", sc.baseURL, sessionId), "application/json", bytes.NewBuffer(b))
//	if err != nil {
//		return fmt.Errorf("failed to make HTTP request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Check if the response status code is 200 OK
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//	return nil
//}

//func (sc *ServiceClient) dbGet(sessionId string, request DBAction) (string, error) {
//
//	b, err := json.Marshal(request)
//	if err != nil {
//		return "", fmt.Errorf("failed to marshal JSON: %w", err)
//	}
//	resp, err := sc.httpClient.Post(fmt.Sprintf("%s/v1/system/%s/db/get", sc.baseURL, sessionId), "application/json", bytes.NewBuffer(b))
//	if err != nil {
//		return "", fmt.Errorf("failed to make HTTP request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Check if the response status code is 200 OK
//	if resp.StatusCode != http.StatusOK {
//		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
//	}
//
//	//return the response body as a string
//	respBody := &bytes.Buffer{}
//	_, err = respBody.ReadFrom(resp.Body)
//	if err != nil {
//		return "", fmt.Errorf("failed to read response body: %w", err)
//	}
//	return respBody.String(), nil
//}
