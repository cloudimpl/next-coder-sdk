package polycode

import (
	"fmt"
	"os"
)

var clientEnv *ClientEnv = nil

func InitClientEnv() {
	if clientEnv != nil {
		return
	}

	appName := os.Getenv("polycode_APP_NAME")
	appPortStr := os.Getenv("polycode_APP_PORT")
	var appPort uint
	if appPortStr == "" {
		appPort = 9998
	} else {
		_, err := fmt.Sscanf(appPortStr, "%d", &appPort)
		if err != nil {
			panic(err)
		}
	}

	clientEnv = &ClientEnv{
		AppName: appName,
		AppPort: appPort,
	}
}

func GetClientEnv() *ClientEnv {
	if clientEnv == nil {
		InitClientEnv()
	}
	return clientEnv
}
