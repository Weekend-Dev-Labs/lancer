package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type LancerConfig struct {
	Database struct {
		Address  string `yaml:"address"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Migrate  bool   `yaml:"migrate"`
	}
	UseRedis bool   `yaml:"use-redis"`
	Redis    string `yaml:"redis"`
	Store    struct {
		Local struct {
			Path string `yaml:"path"`
			Temp string `yaml:"temp-path"`
		}
	}
	AuthEndpoint string `yaml:"auth-endpoint"`
}

func ParseFlags() *LancerConfig {
	cfg := &LancerConfig{}
	filePath := ""

	flag.StringVar(&filePath, "config", "", "Sets the Lancer Configuration")

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

	flag.Parse()

	if filePath != "" {
		// Load the config file

		data, err := os.ReadFile(filePath)

		if err != nil {
			log.Fatalf("[Lancer Error] Failed to load the file (%v)", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("[Lancer Error] Invalid file content for configuration (%v)", err)
		}

		fmt.Println(cfg.Database.Name)
	}

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

	return cfg
}

func (c *LancerConfig) GetDatabaseConnectionString() string {
	return "postgresql://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Address + "/" + c.Database.Name + "?sslmode=disable"
}
