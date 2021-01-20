package service

import (
	"net/mail"
	"net/url"

	"go.uber.org/zap"

	"XustAutoSignIn/model"
	"XustAutoSignIn/serializer"
	"XustAutoSignIn/utils"
	"XustAutoSignIn/variable"
)

type XustCreateUserService struct {
	Url        string `json:"url" binding:"required,url,startswith=http://ehallplatform.xust.edu.cn|startswith=https://ehallplatform.xust.edu.cn"`
	Sno        string `json:"sno" binding:"required,numeric,len=11"`
	Email      string `json:"email" binding:"required,email"`
	Time       int    `json:"time" binding:"required,oneof=1 2"`
	VerifyCode string `json:"verifyCode" binding:"required,numeric,len=6"`
}

func (this *XustCreateUserService) Server() (serializer.ResultUtil, []zap.Field, error) {
	fields := []zap.Field{zap.String("sno", this.Sno)}
	key := utils.GetVerifyCodeKey(this.Sno, this.Email)
	err := checkVerifyCode(key, this.VerifyCode)
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   err.Error(),
			IsSuccess: false,
		}, fields, err
	}
	u, _ := url.Parse(this.Url)
	uid := u.Query()["uid"][0]
	name, err := utils.GetXustUserName(uid, this.Sno)
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedGetUserName,
			IsSuccess: false,
		}, fields, err
	}
	user := &model.XUSTUser{
		Uid:   uid,
		Sno:   this.Sno,
		Email: this.Email,
		Name:  name,
		Time:  this.Time,
	}
	err = variable.DB.Create(&user).Error
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedCreateUser,
			IsSuccess: false,
		}, fields, err
	}
	variable.Cache.Delete(key)
	err = utils.SendEmail(
		mail.Address{Name: this.Sno, Address: this.Email}, "打卡任务已添加",
		"打卡任务添加成功！本系统已启用邮件通知服务，打卡失败会邮件通知，请放心使用！")
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedSendEmail,
			IsSuccess: false,
		}, fields, err
	}
	return serializer.ResultUtil{
		Code:      1000,
		Message:   variable.SuccessCreateUser,
		IsSuccess: true,
	}, fields, nil
}

type XustDeleteUserService struct {
	Sno        string `json:"sno" binding:"required,numeric,len=11"`
	Email      string `json:"email" binding:"required,email"`
	VerifyCode string `json:"verifyCode" binding:"required,numeric,len=6"`
}

func (this *XustDeleteUserService) Server() (serializer.ResultUtil, []zap.Field, error) {
	fields := []zap.Field{zap.String("sno", this.Sno)}
	key := utils.GetVerifyCodeKey(this.Sno, this.Email)
	err := checkVerifyCode(key, this.VerifyCode)
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   err.Error(),
			IsSuccess: false,
		}, fields, err
	}
	err = variable.DB.Where("sno = ?", this.Sno).Delete(&model.XUSTUser{}).Error
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedDeleteUser,
			IsSuccess: false,
		}, fields, err
	}
	variable.Cache.Delete(key)
	err = utils.SendEmail(
		mail.Address{Name: this.Sno, Address: this.Email}, "打卡任务已删除",
		"打卡任务删除成功！有缘再见！")
	if err != nil {
		return serializer.ResultUtil{
			Code:      1001,
			Message:   variable.FailedSendEmail,
			IsSuccess: false,
		}, fields, err
	}
	return serializer.ResultUtil{
		Code:      1000,
		Message:   variable.SuccessDeleteUser,
		IsSuccess: true,
	}, fields, nil
}
