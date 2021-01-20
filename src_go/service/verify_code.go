package service

import (
	"net/mail"

	"go.uber.org/zap"

	"XustAutoSignIn/serializer"
	"XustAutoSignIn/utils"
	"XustAutoSignIn/variable"
)

type VerifyCodeService struct {
	Sno   string `form:"sno" binding:"required,numeric"`
	Email string `form:"email" binding:"required,email"`
}

func (this *VerifyCodeService) Server() (serializer.ResultUtil, []zap.Field, error) {
	fields := []zap.Field{zap.String("email", this.Email)}
	_ = utils.RandomVerifyCode(this.Sno, this.Email)
	verifyCode := utils.RandomVerifyCode(this.Sno, this.Email)
	err := utils.SendEmail(
		mail.Address{Name: this.Sno, Address: this.Email}, "验证码",
		"您的验证码为："+verifyCode)
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedVerifyCode,
			IsSuccess: false,
		}, fields, err
	}
	return serializer.ResultUtil{
		Code:      1000,
		Message:   variable.SuccessVerifyCode,
		IsSuccess: true,
	}, fields, nil
}
