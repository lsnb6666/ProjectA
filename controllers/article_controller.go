package controllers

import (
	"ProjectA/config"
	"ProjectA/global"
	"ProjectA/models"
	"ProjectA/utils"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var articles models.Article
	if err := ctx.ShouldBindJSON(&articles); err != nil {
		config.Logger.Warn("Invalid create article request",
			zap.Error(err),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.BadRequest(ctx, err.Error())
		return
	}

	if err := global.Db.Create(&articles).Error; err != nil {
		config.Logger.Error("Failed to create article",
			zap.Error(err),
			zap.String("title", articles.Title),
		)
		utils.InternalServerError(ctx, "failed to create article")
		return
	}

	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		config.Logger.Warn("Failed to delete cache",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	}

	config.Logger.Info("Article created successfully",
		zap.String("title", articles.Title),
		zap.Uint("id", articles.ID),
	)
	utils.Success(ctx, articles)
}

func GetArticles(ctx *gin.Context) {
	cacheData, err := global.RedisDB.Get(cacheKey).Result()

	if err == redis.Nil {
		config.Logger.Debug("Cache miss - querying database",
			zap.String("cache_key", cacheKey),
		)

		var articles []models.Article
		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.NotFound(ctx, "no articles found")
			} else {
				config.Logger.Error("Failed to query articles", zap.Error(err))
				utils.InternalServerError(ctx, "failed to query articles")
			}
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			config.Logger.Error("Failed to marshal articles", zap.Error(err))
			utils.InternalServerError(ctx, "failed to process data")
			return
		}

		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			config.Logger.Warn("Failed to set cache",
				zap.Error(err),
				zap.String("cache_key", cacheKey),
			)
		}
		utils.Success(ctx, articles)
		return
	} else if err != nil {
		config.Logger.Error("Redis query failed",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
		utils.InternalServerError(ctx, "cache query failed")
		return
	}

	var articles []models.Article
	if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
		config.Logger.Error("Failed to unmarshal cache",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
		utils.InternalServerError(ctx, "failed to process cache")
		return
	}
	config.Logger.Debug("Cache hit - returning articles",
		zap.Int("count", len(articles)),
		zap.String("cache_key", cacheKey),
	)
	utils.Success(ctx, articles)
}

func GetArticlesByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article
	if err := global.Db.Where("id=?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config.Logger.Warn("Article not found",
				zap.String("id", id),
			)
			utils.NotFound(ctx, "article not found")
		} else {
			config.Logger.Error("Failed to query article",
				zap.Error(err),
				zap.String("id", id),
			)
			utils.InternalServerError(ctx, "failed to query article")
		}
		return
	}
	utils.Success(ctx, article)
}
