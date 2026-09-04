package controllers

import (
	"ProjectA/config"
	"ProjectA/global"
	"ProjectA/models"
	"ProjectA/utils"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		config.Logger.Warn("Invalid create exchange rate request",
			zap.Error(err),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.BadRequest(ctx, err.Error())
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		config.Logger.Error("Failed to create exchange rate",
			zap.Error(err),
			zap.String("from", exchangeRate.FromCurrency),
			zap.String("to", exchangeRate.ToCurrency),
		)
		utils.InternalServerError(ctx, "failed to create exchange rate")
		return
	}

	config.Logger.Info("Exchange rate created successfully",
		zap.String("from", exchangeRate.FromCurrency),
		zap.String("to", exchangeRate.ToCurrency),
		zap.Float64("rate", exchangeRate.Rate),
	)
	utils.Success(ctx, exchangeRate)
}

func GetExchangeRates(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate

	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFound(ctx, "no exchange rates found")
		} else {
			config.Logger.Error("Failed to query exchange rates",
				zap.Error(err),
			)
			utils.InternalServerError(ctx, "failed to query exchange rates")
		}
		return
	}

	utils.Success(ctx, exchangeRates)
}
