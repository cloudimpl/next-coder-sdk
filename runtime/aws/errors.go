package aws

import "github.com/CloudImpl-Inc/next-coder-sdk/polycode"

var ErrInvalidEvent = polycode.DefineError("polycode.aws.runtime", 1, "invalid event type ")
var ErrDecodeFailed = polycode.DefineError("polycode.aws.runtime", 2, "invalid msg format")
var ErrRuntimeError = polycode.DefineError("polycode.aws.runtime", 3, "runtime error")
var ErrHttpHandlerNotSet = polycode.DefineError("polycode.aws.runtime", 4, "http handler not set")
