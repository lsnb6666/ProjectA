package controllers

import (
	"ProjectA/config"
	"ProjectA/global"
	"ProjectA/models"
	"ProjectA/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		config.Logger.Warn("Invalid register request",
			zap.Error(err),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.BadRequest(ctx, err.Error())
		return
	}

	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		config.Logger.Error("Failed to hash password",
			zap.Error(err),
			zap.String("username", user.Username),
		)
		utils.InternalServerError(ctx, "password encryption failed")
		return
	}
	user.Password = hashedPwd

	if err := global.Db.Create(&user).Error; err != nil {
		config.Logger.Error("Failed to create user",
			zap.Error(err),
			zap.String("username", user.Username),
		)
		utils.InternalServerError(ctx, "failed to create user")
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		config.Logger.Error("Failed to generate JWT",
			zap.Error(err),
			zap.String("username", user.Username),
		)
		utils.InternalServerError(ctx, "failed to generate token")
		return
	}

	config.Logger.Info("User registered successfully",
		zap.String("username", user.Username),
		zap.String("ip", ctx.ClientIP()),
	)
	utils.Success(ctx, gin.H{"token": token})
}

func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		config.Logger.Warn("Invalid login request",
			zap.Error(err),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.BadRequest(ctx, err.Error())
		return
	}

	var user models.User
	if err := global.Db.Where("username=?", input.Username).First(&user).Error; err != nil {
		config.Logger.Warn("Login failed - user not found",
			zap.String("username", input.Username),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.Unauthorized(ctx, "wrong credentials")
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		config.Logger.Warn("Login failed - wrong password",
			zap.String("username", input.Username),
			zap.String("ip", ctx.ClientIP()),
		)
		utils.Unauthorized(ctx, "wrong credentials")
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		config.Logger.Error("Failed to generate JWT",
			zap.Error(err),
			zap.String("username", user.Username),
		)
		utils.InternalServerError(ctx, "failed to generate token")
		return
	}

	config.Logger.Info("User logged in successfully",
		zap.String("username", user.Username),
		zap.String("ip", ctx.ClientIP()),
	)
	utils.Success(ctx, gin.H{"token": token})
}
