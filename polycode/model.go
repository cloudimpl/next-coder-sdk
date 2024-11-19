package polycode

import (
	"time"
)

type BackoffStrategy struct {
	InitialInterval time.Duration `json:"initialInterval"`
	MaxInterval     time.Duration `json:"maxInterval"`
	Multiplier      float64       `json:"multiplier"`
}

type TaskOptions struct {
	Timeout         time.Duration   `json:"timeout"`
	Retries         int             `json:"retries"`
	RetryOnFail     bool            `json:"retryOnFail"`
	BackoffStrategy BackoffStrategy `json:"backoffStrategy"`
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
	IsError bool     `json:"isError"`
	Output  any      `json:"output"`
	Error   Error    `json:"error"`
	Logs    []LogMsg `json:"logs"`
}

type ApiStartEvent struct {
	SessionId string     `json:"sessionId"`
	Request   ApiRequest `json:"request"`
}

type ApiCompleteEvent struct {
	Path     string      `json:"path"`
	Response ApiResponse `json:"response"`
	Logs     []LogMsg    `json:"logs"`
}

type ErrorEvent struct {
	Error Error `json:"error"`
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
}

type ClientEnv struct {
	EnvId   string `json:"envId"`
	AppName string `json:"appName"`
	AppPort uint   `json:"appPort"`
}
