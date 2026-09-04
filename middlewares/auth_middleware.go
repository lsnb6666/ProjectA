package middlewares

import (
	"ProjectA/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			utils.Unauthorized(ctx, "missing authorization header")
			ctx.Abort()
			return
		}

		username, err := utils.ParseJWT(token)
		if err != nil {
			utils.Unauthorized(ctx, "invalid token")
			ctx.Abort()
			return
		}

		ctx.Set("username", username)
		ctx.Next()
	}
}
