package polycode

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type StartAppRequest struct {
	ClientPort int `json:"clientPort"`
}

type ExecRequest struct {
	ServiceId  string      `json:"serviceId"`
	EntryPoint string      `json:"entryPoint"`
	Options    TaskOptions `json:"options"`
	Input      TaskInput   `json:"input"`
}

// PutRequest represents the JSON structure for put operations
type PutRequest struct {
	Action     string                 `json:"action"`
	Collection string                 `json:"collection"`
	Key        string                 `json:"key"`
	Item       map[string]interface{} `json:"item"`
}

// QueryRequest represents the JSON structure for query operations
type QueryRequest struct {
	Collection string        `json:"collection"`
	Key        string        `json:"key"`
	Filter     string        `json:"filter"`
	Args       []interface{} `json:"args"`
	Limit      int           `json:"limit"`
}

// GetFileRequest represents the JSON structure for get file operations
type GetFileRequest struct {
	Key string `json:"key"`
}

// GetFileResponse represents the JSON structure for get file response
type GetFileResponse struct {
	Content string `json:"content"`
}

// PutFileRequest represents the JSON structure for put file operations
type PutFileRequest struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

// ServiceClient is a reusable client for calling the service API
type ServiceClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewServiceClient creates a new ServiceClient with a reusable HTTP client
func NewServiceClient(baseURL string) *ServiceClient {
	return &ServiceClient{
		httpClient: &http.Client{
			Timeout: time.Second * 30, // Set a reasonable timeout for HTTP requests
		},
		baseURL: baseURL,
	}
}

// StartApp starts the app
func (sc *ServiceClient) StartApp(req StartAppRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, "", "v1/system/app/start", req)
}

// ExecService executes a service with the given request
func (sc *ServiceClient) ExecService(sessionId string, req ExecRequest) (TaskOutput, error) {
	var res TaskOutput
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/service/exec", req, &res)
	if err != nil {
		return TaskOutput{}, err
	}

	if res.IsAsync {
		panic(ErrTaskInProgress)
	}
	return res, nil
}

// GetItem gets an item from the database
func (sc *ServiceClient) GetItem(sessionId string, req QueryRequest) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/get", req, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// QueryItems queries items from the database
func (sc *ServiceClient) QueryItems(sessionId string, req QueryRequest) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/query", req, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PutItem puts an item into the database
func (sc *ServiceClient) PutItem(sessionId string, req PutRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/db/put", req)
}

// GetFile gets a file from the file store
func (sc *ServiceClient) GetFile(sessionId string, req GetFileRequest) (GetFileResponse, error) {
	var res GetFileResponse
	err := executeApiWithResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/get", req, &res)
	if err != nil {
		return GetFileResponse{}, err
	}
	return res, nil
}

// PutFile puts a file into the file store
func (sc *ServiceClient) PutFile(sessionId string, req PutFileRequest) error {
	return executeApiWithoutResponse(sc.httpClient, sc.baseURL, sessionId, "v1/context/file/put", req)
}

func executeApiWithoutResponse(httpClient *http.Client, baseUrl string, sessionId string, path string, req any) error {
	println(fmt.Sprintf("client: exec api without response from %s with session id %s", path, sessionId))

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, path), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-polycode-task-session-id", sessionId)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error, status: %v", resp.Status)
	}

	return nil
}

func executeApiWithResponse[T any](httpClient *http.Client, baseUrl string, sessionId string, path string, req any, res *T) error {
	println(fmt.Sprintf("client: exec api with response from %s with session id %s", path, sessionId))

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, path), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-polycode-task-session-id", sessionId)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error, status: %v", resp.Status)
	}

	if res != nil {
		err = json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			return err
		}
	}

	return nil
}
