package polycode

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

var clientEnv *ClientEnv = nil

func loadIni() {
	cfg, err := ini.Load("env.ini")
	if err != nil {
		println("Failed to load INI file: %v", err)
		return
	}

	// Loop through all sections and keys.
	// Note: You might want to skip the default section if needed.
	for _, section := range cfg.Sections() {
		// Optionally, skip the default section if it doesn't need processing:
		// if section.Name() == ini.DEFAULT_SECTION {
		//     continue
		// }
		for _, key := range section.Keys() {
			// Set each key as an environment variable.
			err := os.Setenv(key.Name(), key.String())
			if err != nil {
				fmt.Printf("Error setting environment variable %s: %v\n", key.Name(), err)
			} else {
				fmt.Printf("Environment variable %s set to %s\n", key.Name(), key.String())
			}
		}
	}

	fmt.Println("Environment variables loaded from INI file.")
}

func initClientEnv() {
	if clientEnv != nil {
		return
	}

	loadIni()

	appName := os.Getenv("polycode_APP_NAME")
	if appName == "" {
		appName = "overridden"
	}

	appPortStr := os.Getenv("polycode_APP_PORT")
	if appPortStr == "" {
		appPortStr = "9998"
	}

	var appPort uint
	_, err := fmt.Sscanf(appPortStr, "%d", &appPort)
	if err != nil {
		appPort = 9998
	}

	clientEnv = &ClientEnv{
		AppName: os.Getenv("polycode_APP_NAME"),
		AppPort: appPort,
	}
}

func GetClientEnv() *ClientEnv {
	if clientEnv == nil {
		initClientEnv()
	}
	return clientEnv
}
