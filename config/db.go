package config

import (
	"ProjectA/global"
	"ProjectA/models"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	Logger.Info("Initializing database connection...")

	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		Logger.Fatal("Failed to initialize database", zap.Error(err), zap.String("dsn", dsn))
	}

	global.Db = db
	Logger.Info("Database connected successfully")

	Logger.Info("Running database migration...")
	if err := global.Db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.ExchangeRate{},
	); err != nil {
		Logger.Fatal("Failed to migrate database", zap.Error(err))
	}
	Logger.Info("Database migration completed successfully")

	sqlDB, err := db.DB()

	if err != nil {
		Logger.Fatal("Failed to configure database", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(AppConfig.Database.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	Logger.Info("Database connection pool configured",
		zap.Int("max_idle_conns", AppConfig.Database.SetMaxIdleConns),
		zap.Int("max_open_conns", AppConfig.Database.SetMaxOpenConns),
	)
}
