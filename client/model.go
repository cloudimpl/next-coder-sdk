package client

import (
	"github.com/CloudImpl-Inc/next-coder-sdk/client/db"
	"github.com/CloudImpl-Inc/next-coder-sdk/polycode"
)

const (
	APIGateway CompletionType = iota
	ServiceCall
)

type CompletionType int

type TaskCompleteEvent struct {
	//Tx             *TxRequest
	//CompletionType CompletionType
	Output         polycode.TaskOutput
	Tx             *db.Tx
	TaskInProgress bool
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type TaskContext struct {
	Id         string `json:"id"`
	SessionId  string `json:"sessionId"`
	EntryPoint string `json:"entryPoint"`
}

type StartAppRequest struct {
}
