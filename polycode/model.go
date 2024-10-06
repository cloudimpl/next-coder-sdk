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

type TaskOutput struct {
	IsAsync bool   `json:"isAsync"`
	IsNull  bool   `json:"isNull"`
	Output  any    `json:"output"`
	Error   *Error `json:"error"`
}
