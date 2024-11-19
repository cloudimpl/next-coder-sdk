package polycode

import (
	"encoding/json"
	"fmt"
	"os"
)

func runCliCommand(args []string) error {
	switch args[0] {
	case "info":
		{
			return getAppInfo(args[1])
		}
	default:
		{
			return fmt.Errorf("unknown command: %s", args[0])
		}
	}
}

func getAppInfo(filePath string) error {
	var services []ServiceData
	for srvName, srv := range serviceMap {
		serviceData := ServiceData{
			Name:  srvName,
			Tasks: make([]TaskData, 0),
		}

		res, err := srv.ExecuteService(nil, "@definition", nil)
		if err != nil {
			return err
		}

		taskList := res.([]string)
		for _, taskName := range taskList {
			taskData := TaskData{
				Name:       taskName,
				IsWorkflow: srv.IsWorkflow(taskName),
			}
			serviceData.Tasks = append(serviceData.Tasks, taskData)
		}

		services = append(services, serviceData)
	}

	req := StartAppRequest{
		Services: services,
		Routes:   loadRoutes(),
	}

	output, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		return err
	}

	return nil
}
