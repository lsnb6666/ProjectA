package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn             string
		SetMaxIdleConns int
		SetMaxOpenConns int
	}
}

var AppConfig *Config

func InitConfig() {
	InitLogger()
	Logger.Info("Starting application...")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		Logger.Fatal("Error reading config file", zap.Error(err))
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		Logger.Fatal("Unable to decode into struct", zap.Error(err))
	}
	Logger.Info("Config loaded successfully",
		zap.String("app_name", AppConfig.App.Name),
		zap.String("port", AppConfig.App.Port),
	)
	initDB()
	initRedis()
}
