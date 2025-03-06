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
		return
	}

	fmt.Println("[LANCER INFO] Creating default configuration here: ", dir)
	err = os.MkdirAll(dir, 0755)

	if err != nil {
		log.Fatalf("[LANCER ERROR] Failed to create configuration")
	}

	defaultYaml := `
port: "8080"
allow-origin:
  - ""
database:
  migrate: false
  address: "localhost:5432"
  user: "postgres"
  password: "qwertyuiop"
  name: "lancer-main"
use-redis: false
redis: ""
server-auth: ""
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
    clientId: ""
    clientSecret: ""
admin-token-secret: ""
auth:
  email: "lancer@email.com"
  password: "password"
`

	if err := os.WriteFile(fmt.Sprintf("%s/lancer.yaml", dir), []byte(defaultYaml), 0755); err != nil {
		log.Fatalf("[LANCER ERROR] Failed to create default configuration yaml file")
	}

	if err := os.WriteFile(filepath.Join(dir, "secrets"), []byte(""), 0755); err != nil {
		fmt.Println("[LANCER ERROR] Failed to create secrets")
	}

	if err := os.WriteFile(filepath.Join(dir, "history"), []byte(""), 0755); err != nil {
		fmt.Println("[LANCER ERROR] Failed to history secrets")
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

func getContent(name string) string {

	filePath := filepath.Join(getPlatfromConfigDirectory(), name)

	content, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("[LANCER WARNING] No %s file found\n", name)
		return ""
	}

	contentString := string(content)

	return contentString
}

func writeContent(content string, name string) {
	os.WriteFile(filepath.Join(getPlatfromConfigDirectory(), name), []byte(content), 0755)
}

func GetHistoryContent() string {
	return getContent(types.AppHistory)
}

func GetSecretContent() string {
	return getContent(types.AppSecrets)
}
