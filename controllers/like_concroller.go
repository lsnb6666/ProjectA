package controllers

import (
	"ProjectA/config"
	"ProjectA/global"
	"ProjectA/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		config.Logger.Error("Failed to like article",
			zap.Error(err),
			zap.String("article_id", articleID),
		)
		utils.InternalServerError(ctx, "failed to like article")
		return
	}

	config.Logger.Info("Article liked successfully",
		zap.String("article_id", articleID),
		zap.String("ip", ctx.ClientIP()),
	)
	utils.SuccessWithMessage(ctx, "Successfully liked the article", nil)
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"

	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		config.Logger.Error("Failed to get like count",
			zap.Error(err),
			zap.String("article_id", articleID),
		)
		utils.InternalServerError(ctx, "failed to get like count")
		return
	}

	utils.Success(ctx, gin.H{"likes": likes})
}
