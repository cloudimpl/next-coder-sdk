package polycode

import (
	"fmt"
	"log"
	"os"
)

var clientEnv *ClientEnv = nil

func InitClientEnv() {
	if clientEnv != nil {
		return
	}

	appPortStr := os.Getenv("polycode_APP_PORT")
	if appPortStr == "" {
		appPortStr = "9998"
	}

	var appPort uint
	_, err := fmt.Sscanf(appPortStr, "%d", &appPort)
	if err != nil {
		log.Fatalf("invalid APP_PORT: %s", appPortStr)
	}

	clientEnv = &ClientEnv{
		EnvId:   os.Getenv("polycode_ENV_ID"),
		AppName: os.Getenv("polycode_APP_NAME"),
		AppPort: appPort,
	}
}

func GetClientEnv() *ClientEnv {
	if clientEnv == nil {
		InitClientEnv()
	}
	return clientEnv
}
