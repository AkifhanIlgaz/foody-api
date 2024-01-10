package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MongoUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	MongoDbName string `mapstructure:"MONGO_DATABASE_NAME"`

	RedisUrl string `mapstructure:"REDIS_URL"`

	Port string `mapstructure:"PORT"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &config, nil
}
