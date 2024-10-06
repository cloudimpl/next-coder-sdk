package client

import "github.com/CloudImpl-Inc/next-coder-sdk/polycode"

var ErrTaskNotFound = polycode.DefineError("polycode.client", 1, "task not found .task: [%s]")
var ErrBadRequest = polycode.DefineError("polycode.client", 2, "bad request")
var ErrTaskExecError = polycode.DefineError("polycode.client", 3, "task execution error")
var ErrTaskInProgress = polycode.DefineError("polycode.client", 4, "task in progress")
