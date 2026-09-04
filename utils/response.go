package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(ctx *gin.Context, message string) {
	Error(ctx, http.StatusBadRequest, message)
}

func Unauthorized(ctx *gin.Context, message string) {
	Error(ctx, http.StatusUnauthorized, message)
}

func Forbidden(ctx *gin.Context, message string) {
	Error(ctx, http.StatusForbidden, message)
}

func NotFound(ctx *gin.Context, message string) {
	Error(ctx, http.StatusNotFound, message)
}

func InternalServerError(ctx *gin.Context, message string) {
	Error(ctx, http.StatusInternalServerError, message)
}
