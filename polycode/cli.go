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
	services, err := ExtractServiceDescription()
	if err != nil {
		return err
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
