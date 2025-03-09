package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/weekend-dev-labs/lancer/types"
	"github.com/weekend-dev-labs/lancer/utils"
	"gopkg.in/yaml.v3"
)

const Version = "v2.0.4"

type LancerConfig struct {
	Port        string   `yaml:"port"`
	AllowOrigin []string `yaml:"allow-origin"`
	Database    struct {
		Address  string `yaml:"address"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Migrate  bool
	}
	UseRedis        bool   `yaml:"use-redis"`
	Redis           string `yaml:"redis"`
	WebhookEndpoint string `yaml:"webhook-endpoint"`
	Store           struct {
		Local struct {
			Path string `yaml:"path"`
			Temp string `yaml:"temp-path"`
		}
		AWS struct {
			Store        bool   `yaml:"store"`
			Bucket       string `yaml:"bucket"`
			Region       string `yaml:"region"`
			ClientID     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
		}
	}
	ServerAuth   bool   `yaml:"server-auth"`
	AuthEndpoint string `yaml:"auth-endpoint"`

	AdminTokenSigningSecret string `yaml:"admin-token-secret"`

	Auth struct {
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	}

	MetricsID uuid.UUID

	IsAwsProvided        bool
	WebhookSigningSecret string
}

func ParseFlags() *LancerConfig {

	cfg := &LancerConfig{}
	filePath := ""

	flag.StringVar(&filePath, "config", "", "Sets the Lancer Configuration")

	// Server Config
	flag.StringVar(&cfg.Port, "port", "8080", "Sets the port in which lancer will run (default 8080)")

	// Database CLI Args
	flag.BoolVar(&cfg.Database.Migrate, "database-migrate", false, "Sets the option for migration")
	flag.StringVar(&cfg.Database.Address, "database-addr", "localhost:5432", "Sets the database address.")
	flag.StringVar(&cfg.Database.User, "database-user", "", "Sets the username for the database connection.")
	flag.StringVar(&cfg.Database.Name, "database-name", "lancer-db", "Sets the database name.")
	flag.StringVar(&cfg.Database.Password, "database-password", "qwertyuiop", "Sets the password for database connection.")

	// Redis Args
	flag.BoolVar(&cfg.UseRedis, "redis", false, "Whether to use redis or not.")
	flag.StringVar(&cfg.Redis, "redis-addr", "", "Sets the address for redis.")

	// Store Args
	flag.StringVar(&cfg.Store.Local.Path, "store-local-path", "store", "Sets the path to store the media files locally.")
	flag.StringVar(&cfg.Store.Local.Temp, "store-local-temp", "temp", "Sets the path to store the media files temporariliy")
	flag.StringVar(&cfg.AuthEndpoint, "auth-endpoint", "", "Sets the path for authentication")

	flag.BoolVar(&cfg.ServerAuth, "server-auth", false, "Whether to authenticate or not with the server (Default: false)")

	flag.StringVar(&cfg.WebhookEndpoint, "webhook-endpoint", "", "Sets the path for webhook endpoint.")
	flag.StringVar(&cfg.AdminTokenSigningSecret, "admin-token-secret", "admin-token", "Sets the token signing secret for the admin token")

	flag.StringVar(&cfg.Auth.Email, "email", "lancer@email.com", "Email to login to dashboard")
	flag.StringVar(&cfg.Auth.Password, "password", "password", "Password to login to dashboard")

	flag.BoolVar(&cfg.Store.AWS.Store, "aws-store", false, "Whether to store media files in AWS S3.")
	flag.StringVar(&cfg.Store.AWS.Bucket, "aws-bucket", "", "S3 Bucket name to store file in.")
	flag.StringVar(&cfg.Store.AWS.Region, "aws-region", "", "AWS Region to store file.")
	flag.StringVar(&cfg.Store.AWS.ClientID, "aws-client-id", "", "AWS Client ID.")
	flag.StringVar(&cfg.Store.AWS.ClientSecret, "aws-client-secret", "", "AWS Client secret.")

	flag.Parse()

	configurationPath := getPlatformConfigFilePath()

	var configPath string

	if configurationPath != "" {
		configPath = configurationPath
	}

	if filePath != "" {
		configPath = filePath
	}

	if filePath != "" || configurationPath != "" {

		data, err := os.ReadFile(configPath)

		if err != nil {
			log.Fatalf("[Lancer Error] Failed to load the file (%v)", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("[Lancer Error] Invalid file content for configuration (%v)", err)
		}
	}

	cfg.IsAwsProvided = false

	// Database Related
	if cfg.Database.Address == "" {
		log.Fatalf("[Lancer Error] Database Address can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -database-addr=<database-address>)")
	}

	if cfg.Database.Name == "" {
		log.Fatalf("[Lancer Error] Database Name can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -database-name=<database-name>)")
	}

	if cfg.Database.User == "" {
		log.Fatalf("[Lancer Error] Database User can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -database-user=<database-user>)")
	}

	if cfg.Database.Password == "" {
		log.Fatalf("[Lancer Error] Database Password can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -database-password=<database-password>)")
	}

	// Redis Related
	if cfg.UseRedis {
		if cfg.Redis == "" {
			log.Fatalf("[Lancer Error] Redis Address can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -redis-addr=<redis-address>)")
		}
	}

	// Local path storage related
	if cfg.Store.Local.Path == "" {
		log.Fatalf("[Lancer Error] Local Store path can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -store-local-path=<stores-path-locally>)")
	}

	if cfg.Store.Local.Temp == "" {
		log.Fatalf("[Lancer Error] Local Store temp can't be empty. If using a config file see this reference : (https://lancer.dev/cfg-file) (use -store-local-temp=<temp-path>)")
	}

	if cfg.Store.AWS.Store {
		if cfg.Store.AWS.Bucket == "" {
			log.Fatalf("[Lancer Error] AWS Bucket name can't be empty, If using a config file see this reference : (https://lancer.dev/cfg-file) (use -aws-bucket=<bucket-name>)")
		}

		if cfg.Store.AWS.Region == "" {
			log.Fatalf("[Lancer Error] AWS Region  can't be empty, If using a config file see this reference : (https://lancer.dev/cfg-file) (use -aws-region=<region>)")
		}
	}

	cfg.HandleStandaloneArgs()

	return cfg
}

func (c *LancerConfig) GetDatabaseConnectionString() string {
	return "postgresql://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Address + "/" + c.Database.Name + "?sslmode=disable"
}

func (c *LancerConfig) GetSigningSecret() string {
	currentEndpoints := c.AuthEndpoint + ":" + c.WebhookEndpoint

	content := GetHistoryContent()

	if currentEndpoints != content {
		secret, err := utils.GenerateSecret(50)

		if err != nil {
			log.Fatalf("[LANCER ERROR] Failed to generate signing secret")
		}

		writeContent(secret, types.AppSecrets)
		writeContent(currentEndpoints, types.AppHistory)

		c.WebhookSigningSecret = secret

		return secret
	}

	c.WebhookSigningSecret = getContent(types.AppSecrets)

	return c.WebhookSigningSecret
}

func (c *LancerConfig) HandleStandaloneArgs() {
	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "migrate":
			c.Database.Migrate = true
			return

		case "version":
			fmt.Println(Version)
			os.Exit(1)
		}
	}
}
