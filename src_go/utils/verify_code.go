package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"

	"github.com/patrickmn/go-cache"

	"XustAutoSignIn/variable"
)

func GetVerifyCodeKey(sno, email string) string {
	return hex.EncodeToString(md5.New().Sum([]byte(sno + email)))
}

func RandomVerifyCode(sno, email string) string {
	var buf bytes.Buffer
	for i := 0; i < 6; i++ {
		buf.WriteString(strconv.Itoa(rand.Int() % 10))
	}
	variable.Cache.Set(GetVerifyCodeKey(sno, email), buf.String(), cache.DefaultExpiration)
	return buf.String()
}
