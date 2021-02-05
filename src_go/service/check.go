package service

import (
	"errors"

	"XustAutoSignIn/variable"
)

func checkVerifyCode(key, verifyCode string) error {
	v, ok := variable.Cache.Get(key)
	if ok == false {
		return errors.New(variable.VerifyCodeExpired)
	}
	if v != verifyCode {
		return errors.New(variable.VerifyCodeIncorrect)
	}
	return nil
}
