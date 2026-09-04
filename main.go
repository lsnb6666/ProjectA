package main

import (
	"ProjectA/config"
	"ProjectA/router"

	"go.uber.org/zap"
)

func main() {
	config.InitConfig()

	defer config.SyncLogger()

	r := router.SetupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}
	config.Logger.Info("Starting HTTP server",
		zap.String("port", port),
		zap.String("app_name", config.AppConfig.App.Name),
	)
	if err := r.Run(port); err != nil {
		config.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
