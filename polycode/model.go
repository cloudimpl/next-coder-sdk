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
	Timeout         time.Duration    // Maximum time allowed for the task to complete
	Retries         int              // Number of times to retry the task upon failure
	RetryOnFail     bool             // Whether to retry the task automatically on failure
	BackoffStrategy *BackoffStrategy // Backoff strategy for handling retries
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

type TaskInput struct {
	Checksum  uint64 `json:"checksum"`
	NoArg     bool   `json:"noArg"`
	TargetReq string `json:"targetReq"`
}

type TaskStartEvent struct {
	Id           string    `json:"id"`
	SessionId    string    `json:"sessionId"`
	TenantId     string    `json:"tenantId"`
	ServiceName  string    `json:"serviceName"`
	PartitionKey string    `json:"partitionKey"`
	EntryPoint   string    `json:"entryPoint"`
	Input        TaskInput `json:"input"`
}

type ApiRequest struct {
	Id     string            `json:"id"`
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Query  map[string]string `json:"query"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type ApiStartEvent struct {
	SessionId  string     `json:"sessionId"`
	Controller string     `json:"controller"`
	Request    ApiRequest `json:"request"`
}

type TaskOutput struct {
	IsAsync bool  `json:"isAsync"`
	IsNull  bool  `json:"isNull"`
	IsError bool  `json:"isError"`
	Output  any   `json:"output"`
	Error   Error `json:"error"`
}

type TaskCompleteEvent struct {
	Output TaskOutput
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
