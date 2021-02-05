package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"XustAutoSignIn/serializer"
	"XustAutoSignIn/service"
	"XustAutoSignIn/variable"
)

var paramError = serializer.ResultUtil{
	Code:      1001,
	Message:   variable.ParamError,
	IsSuccess: false,
}

func Ping(ctx *gin.Context) {
	variable.Logger.
		With(zap.String("module", "router"), zap.String("method", "GET"), zap.String("uri", "/api/ping"), zap.String("Source", ctx.ClientIP())).
		Info("Ping, Pong, Peng")
	ctx.String(http.StatusOK, "Pong")
}

// 发送验证码
func VerificationCode(ctx *gin.Context) {
	logger := variable.Logger.With(zap.String("module", "router"), zap.String("method", "GET"), zap.String("uri", "/api/verifyCode"), zap.String("Source", ctx.ClientIP()))
	var s service.VerifyCodeService
	if err := ctx.ShouldBindQuery(&s); err != nil {
		logger.Warn(err.Error())
		ctx.JSON(http.StatusBadRequest, paramError)
		return
	}
	data, fields, err := s.Server()
	if err != nil {
		logger.Desugar().Error(err.Error(), fields...)
	} else {
		logger.Desugar().Info(variable.SuccessVerifyCode, fields...)
	}
	ctx.JSON(http.StatusOK, data)
}

// 新建用户
func XUSTCreateUser(ctx *gin.Context) {
	logger := variable.Logger.
		With(zap.String("module", "router"), zap.String("method", "PUT"), zap.String("uri", "/api/xust/user"), zap.String("Source", ctx.ClientIP()))
	var s service.XustCreateUserService
	if err := ctx.ShouldBindJSON(&s); err != nil {
		logger.Warn(err.Error())
		ctx.JSON(http.StatusBadRequest, paramError)
		return
	}
	data, fields, err := s.Server()
	if err != nil {
		logger.Desugar().Error(err.Error(), fields...)
	} else {
		logger.Desugar().Info(variable.SuccessCreateUser, fields...)
	}
	ctx.JSON(http.StatusOK, data)
}

// 删除用户
func XUSTDeleteUser(ctx *gin.Context) {
	logger := variable.Logger.
		With(zap.String("module", "router"), zap.String("method", "DELETE"), zap.String("uri", "/api/xust/user"), zap.String("Source", ctx.ClientIP()))
	var s service.XustDeleteUserService
	if err := ctx.ShouldBindJSON(&s); err != nil {
		logger.Warn(err.Error())
		ctx.JSON(http.StatusBadRequest, paramError)
		return
	}
	data, fields, err := s.Server()
	if err != nil {
		logger.Desugar().Error(err.Error(), fields...)
	} else {
		logger.Desugar().Info(variable.SuccessDeleteUser, fields...)
	}
	ctx.JSON(http.StatusOK, data)
}
