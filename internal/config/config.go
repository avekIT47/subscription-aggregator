package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Port     string `yaml:"port"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(path string) *Config {
	log.Printf("Loading config from %s", path)

	cfg := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't read config.yaml: %v", err)
	}
	log.Printf("Config file %s successfully read, size: %d bytes", path, len(file))

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		log.Fatalf("couldn't parse config.yaml: %v", err)
	}
	log.Printf("Config file %s successfully parsed", path)

	return cfg
}

func (context *Config) GetDSN() string {
	db := context.Database
	dsn := "host=" + db.Host +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" port=" + db.Port +
		" sslmode=" + db.SSLMode

	log.Printf("Generated DSN for DB connection: %s", dsn)
	return dsn
}
