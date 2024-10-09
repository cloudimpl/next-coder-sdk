package polycode

import "fmt"

var currentRuntime Runtime = nil

//func GetRuntimeInfo() RuntimeInfo {
//	if currentRuntime == nil {
//		panic("Current runtime not set")
//	}
//	return currentRuntime
//}

//type RuntimeInfo interface {
//	Name() string
//}

type Runtime interface {
	Name() string
	AppConfig() AppConfig
	Start(params []any) error
}

// SetCurrentRuntime Set the current runtime: this is used inside .polycode/runtime.go
func SetCurrentRuntime(runtime Runtime) {
	if currentRuntime != nil {
		println("Current runtime already set.ignored")
	}
	fmt.Printf("runtime %s set\n", runtime.Name())
	currentRuntime = runtime
}

// Start the runtime: this is used in client main function
func Start(params ...any) error {
	if currentRuntime == nil {
		println("runtime not set")
		return fmt.Errorf("runtime not set")
	}
	return currentRuntime.Start(params)
}
