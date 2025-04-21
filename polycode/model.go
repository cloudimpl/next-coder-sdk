package polycode

import (
	"time"
)

type DescribeMethodRequest struct {
	Method string
}

type DescribeMethodResponse struct {
	Method     string      `json:"method"`
	IsWorkflow bool        `json:"isWorkflow"`
	Input      interface{} `json:"input"`
}

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

type MethodStartEvent struct {
	SessionId string      `json:"sessionId"`
	Method    string      `json:"method"`
	Meta      ContextMeta `json:"meta"`
	Input     any         `json:"input"`
}

type ServiceStartEvent struct {
	SessionId string      `json:"sessionId"`
	Service   string      `json:"service"`
	Method    string      `json:"method"`
	Meta      ContextMeta `json:"meta"`
	Input     any         `json:"input"`
}

type ServiceMeta struct {
	IsWorkflow bool `json:"isWorkflow"`
}

type ServiceCompleteEvent struct {
	IsError bool        `json:"isError"`
	Output  any         `json:"output"`
	Error   Error       `json:"error"`
	Logs    []LogMsg    `json:"logs"`
	Meta    ServiceMeta `json:"meta"`
}

type ApiStartEvent struct {
	SessionId string      `json:"sessionId"`
	Meta      ContextMeta `json:"meta"`
	Request   ApiRequest  `json:"request"`
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

type ServiceDescription struct {
	Name  string              `json:"name"`
	Tasks []MethodDescription `json:"tasks"`
}

type MethodDescription struct {
	Name       string      `json:"name"`
	IsWorkflow bool        `json:"isWorkflow"`
	Input      interface{} `json:"input"`
}

type ClientEnv struct {
	EnvId   string `json:"envId"`
	AppName string `json:"appName"`
	AppPort uint   `json:"appPort"`
}

type ContextMeta struct {
	OrgId        string `json:"orgId"`
	EnvId        string `json:"envId"`
	AppName      string `json:"appName"`
	AppId        string `json:"appId"`
	TenantId     string `json:"tenantId"`
	PartitionKey string `json:"partitionKey"`
	TaskGroup    string `json:"taskGroup"`
	TaskName     string `json:"taskName"`
	TaskId       string `json:"taskId"`
	TraceId      string `json:"traceId"`
	InputId      string `json:"inputId"`
}
