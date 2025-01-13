package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/weekend-dev-labs/lancer/types"
)

func InitConfig() {
	dir := getPlatfromConfigDirectory()

	_, err := os.Stat(dir)

	if err == nil {
		fmt.Println("[LANCER INFO] Configuration found loading from : ", dir)
		return
	}

	fmt.Println("[LANCER INFO] Creating default configuration here: ", dir)
	err = os.MkdirAll(dir, 0755)

	if err != nil {
		log.Fatalf("[LANCER ERROR] Failed to create configuration")
	}

	defaultYaml := `
database:
  migrate: false
  address: "localhost:5432"
  user: "postgres"
  password: "qwertyuiop"
  name: "lancer-main"
use-redis: false
redis: ""
webhook-endpoint: ""
auth-endpoint: ""
store:
  local:
    path: "store"
    temp-path: "temp"
  aws:
    store: false
    bucket: ""
    region: ""
    config: ""
admin-token-secret: ""
auth:
  email: "lancer@email.com"
  password: "password"
`

	if err := os.WriteFile(fmt.Sprintf("%s/lancer.yaml", dir), []byte(defaultYaml), 0755); err != nil {
		log.Fatalf("[LANCER ERROR] Failed to create default configuration yaml file")
	}

}

func getPlatfromConfigDirectory() string {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}

	case "darwin":
		baseDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")

	case "linux":
		baseDir = filepath.Join(os.Getenv("HOME"), ".config")
	default:
		baseDir = "."
	}

	return filepath.Join(baseDir, types.AppName)
}

func getPlatformConfigFilePath() string {
	return filepath.Join(getPlatfromConfigDirectory(), types.AppConfigFile)
}
