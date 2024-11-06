package polycode

import (
	"time"
)

type BackoffStrategy struct {
	InitialInterval time.Duration // Initial retry interval
	MaxInterval     time.Duration // Maximum interval between retries
	Multiplier      float64       // Multiplier to apply to the interval after each failure
}

type TaskOptions struct {
	Timeout         time.Duration   // Maximum time allowed for the task to complete
	Retries         int             // Number of times to retry the task upon failure
	RetryOnFail     bool            // Whether to retry the task automatically on failure
	BackoffStrategy BackoffStrategy // Backoff strategy for handling retries
	PartitionKey    string
	TenantId        string
}

func (t TaskOptions) WithPartitionKey(partitionKey string) TaskOptions {
	t.PartitionKey = partitionKey
	return t
}

func (t TaskOptions) WithTimeout(timeout time.Duration) TaskOptions {
	t.Timeout = timeout
	return t
}

type ServiceStartEvent struct {
	SessionId string `json:"sessionId"`
	Service   string `json:"service"`
	Method    string `json:"method"`
	Input     any    `json:"input"`
}

type ServiceCompleteEvent struct {
	IsError bool  `json:"isError"`
	Output  any   `json:"output"`
	Error   Error `json:"error"`
}

type ApiStartEvent struct {
	SessionId string     `json:"sessionId"`
	Request   ApiRequest `json:"request"`
}

type ApiCompleteEvent struct {
	Path     string      `json:"path"`
	Response ApiResponse `json:"response"`
}

type ErrorEvent struct {
	Error Error
}

type ApiRequest struct {
	Id              string            `json:"id"`
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	Query           map[string]string `json:"query"`
	Header          map[string]string `json:"header"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type ApiResponse struct {
	StatusCode      int               `json:"statusCode"`
	Header          map[string]string `json:"header"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type RouteData struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type ServiceData struct {
	Name  string     `json:"name"`
	Tasks []TaskData `json:"tasks"`
}

type TaskData struct {
	Name       string `json:"name"`
	IsWorkflow bool   `json:"isWorkflow"`
	IsReadOnly bool   `json:"isReadOnly"`
}

type ClientEnv struct {
	AppName string `json:"appName"`
	AppPort uint   `json:"appPort"`
}
