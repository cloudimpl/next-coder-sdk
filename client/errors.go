package client

import "github.com/CloudImpl-Inc/next-coder-sdk/polycode"

var ErrBadRequest = polycode.DefineError("polycode.client", 2, "bad request")
var ErrTaskExecError = polycode.DefineError("polycode.client", 3, "task execution error")
var ErrUnknownError = polycode.DefineError("polycode.client", 4, "unknown error")
var ErrPanic = polycode.DefineError("polycode.client", 5, "task in progress")
var ErrTaskInProgress = &ErrPanic
