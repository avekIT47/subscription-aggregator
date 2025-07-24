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
}

func LoadConfig(path string) *Config {
	cfg := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't read config.yaml: %v", err)
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		log.Fatalf("couldn't parse config.yaml: %v", err)
	}

	return cfg
}

func (context *Config) GetDSN() string {
	db := context.Database
	return "host=" + db.Host +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" port=" + db.Port +
		" sslmode=" + db.SSLMode
}
