package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"sync"
	"team-task/pkg/logger"
)

type (
	Config struct {
		Server ServerConfig
		Auth   AuthConfig
		STAN   STANConfig
	}

	ServerConfig struct {
		PortOne string `env:"PORT_ONE" env-default:"8080"`
		PortTwo string `env:"PORT_TWO" env-default:"8081"`
	}

	AuthConfig struct {
		Username string `env:"USERNAME" env-default:"admin"`
		Password string `env:"PASSWORD" env-default:"admin"`
	}

	STANConfig struct {
		ClusterID string `env:"CLUSTER_ID" env-default:"test-cluster"`
		ClientID  string `env:"CLIENT_ID" env-default:"client-id"`
		URL       string `env:"NATS_URL" env-default:"0.0.0.0:4222"`
		Subject   string `env:"NATS_SUBJECT" env-default:"wildberries"`
	}
)

var instance *Config

func GetConfig() *Config {
	once := sync.Once{}
	once.Do(func() {
		instance = &Config{}

		if err := godotenv.Load(".env"); err != nil {
			logger.Errorf("error during loading environment variables: %s\n", err.Error())
			return
		}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Errorf("error during mapping environment variables: %s\n", help)
			return
		}
	})

	return instance
}
