package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	AppEnv      string `mapstructure:"APP_ENV"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      int    `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	RedisHost   string `mapstructure:"REDIS_HOST"`
	RedisPort   int    `mapstructure:"REDIS_PORT"`
	RedisPass   string `mapstructure:"REDIS_PASSWORD"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	JWTExpiry   string `mapstructure:"JWT_EXPIRY"`
	OllamaURL   string `mapstructure:"OLLAMA_URL"`
	OllamaModel string `mapstructure:"OLLAMA_MODEL"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{
		ServerPort:  "8080",
		AppEnv:      "development",
		DBHost:      "localhost",
		DBPort:      3306,
		DBUser:      "app",
		DBPassword:  "apppass",
		DBName:      "visual_mesin",
		RedisHost:   "localhost",
		RedisPort:   6379,
		JWTSecret:   "change-me",
		JWTExpiry:   "24h",
		OllamaURL:   "http://localhost:11434",
		OllamaModel: "llama3",
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
