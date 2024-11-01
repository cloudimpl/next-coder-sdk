package polycode

import (
	"encoding/json"
	"fmt"
)

func runCliCommand(args []string) error {
	switch args[0] {
	case "info":
		{
			return getAppInfo()
		}
	default:
		{
			return fmt.Errorf("unknown command: %s", args[0])
		}
	}
}

func getAppInfo() error {
	var services []ServiceData
	for name := range serviceMap {
		services = append(services, ServiceData{
			Name: name,
			// ToDo: Add task info
		})
	}

	req := StartAppRequest{
		Services: services,
		Routes:   loadRoutes(),
	}

	output, err := json.Marshal(req)
	if err != nil {
		return err
	}

	println(string(output))
	return nil
}
