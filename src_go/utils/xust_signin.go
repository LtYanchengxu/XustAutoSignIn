package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"XustAutoSignIn/model"
	"XustAutoSignIn/variable"

	"go.uber.org/zap"
)

const (
	XUST_TIME_FORMAT    = "2006-01-02 15:04:05"
	XUST_HOST           = "ehallplatform.xust.edu.cn"
	XUST_GET_COOKIE_URL = "https://ehallplatform.xust.edu.cn/default/jkdk/mobile/mobJkdkAdd_test.jsp?uid="
	XUST_GET_USER_INFO  = "http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.getJkdkRownum.biz.ext?gh="
	XUST_SIGNIN         = "http://ehallplatform.xust.edu.cn/default/jkdk/mobile/com.primeton.eos.jkdk.xkdjkdkbiz.jt.biz.ext"
)

const (
	XUST_AFTERNOON = iota
	XUST_NIGHT
)

func getXustCookie(uid string, userAgentId int) (string, error) {
	req, err := http.NewRequest("GET", XUST_GET_COOKIE_URL+uid, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Connection", "close")
	req.Header.Add("User-Agent", variable.UserAgents[userAgentId])
	req.Header.Add("Host", XUST_HOST)
	resp, err := (&http.Client{Transport: &http.Transport{DisableKeepAlives: true}}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return strings.Split(resp.Header.Get("Set-Cookie"), ";")[0], nil
}

func getXustSigninInfo(uid, sno, cookie string, userAgentId int) (map[string]json.RawMessage, error) {
	req, err := http.NewRequest("GET", XUST_GET_USER_INFO+sno, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Connection", "close")
	req.Header.Add("User-Agent", variable.UserAgents[userAgentId])
	req.Header.Add("Host", XUST_HOST)
	req.Header.Add("Cookie", cookie)
	resp, err := (&http.Client{Transport: &http.Transport{DisableKeepAlives: true}}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var obj map[string]json.RawMessage
	if err = json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	var list []json.RawMessage
	_ = json.Unmarshal(obj["list"], &list)
	obj = make(map[string]json.RawMessage)
	_ = json.Unmarshal(list[0], &obj)
	return obj, nil
}

func getXustSignInRequestDate(info map[string]json.RawMessage, tag int) []byte {
	var tmpString string
	var tmpBytes []byte
	nowTime := time.Now().Local()
	tmpString = nowTime.Format(XUST_TIME_FORMAT)
	tmpBytes, _ = json.Marshal(&tmpString)
	info["TBSJ"] = json.RawMessage(tmpBytes)
	tmpString = tmpString[:10]
	tmpBytes, _ = json.Marshal(&tmpString)
	info["TIME"] = json.RawMessage(tmpBytes)
	if tag == XUST_NIGHT {
		nowTime = nowTime.Add(24 * time.Hour)
		tmpString = nowTime.Format(XUST_TIME_FORMAT)[:10]
		tmpBytes, _ = json.Marshal(&tmpString)
	}
	info["JRRQ1"] = json.RawMessage(tmpBytes)
	if tag == XUST_AFTERNOON {
		tmpString = "1"
	} else if tag == XUST_NIGHT {
		tmpString = "0"
	}
	tmpBytes, _ = json.Marshal(&tmpString)
	info["JDLX"] = json.RawMessage(tmpBytes)
	tmpString = "中国"
	tmpBytes, _ = json.Marshal(&tmpString)
	info["GUO"] = json.RawMessage(tmpBytes)
	info["XXDZ41"] = info["XXDZ4_1"]
	delete(info, "XXDZ4_1")
	delete(info, "ID")
	delete(info, "PROCINSTID")
	obj := make(map[string]json.RawMessage)
	for k, v := range info {
		obj[strings.ToLower(k)] = v
	}
	tmpBytes, _ = json.Marshal(&obj)
	post := make(map[string]json.RawMessage)
	post["xkdjkdk"] = json.RawMessage(tmpBytes)
	tmpBytes, _ = json.Marshal(post)
	return tmpBytes
}

func GetXustUserName(uid, sno string) (name string, err error) {
	userAgentId := rand.Int() % len(variable.UserAgents)
	cookie, err := getXustCookie(uid, userAgentId)
	if err != nil {
		return
	}
	obj, err := getXustSigninInfo(uid, sno, cookie, userAgentId)
	if err != nil {
		return
	}
	_ = json.Unmarshal(obj["XM"], &name)
	return
}

func xustSignIn(uid, cookie string, userAgentId int, post []byte) error {
	req, err := http.NewRequest("POST", XUST_SIGNIN, bytes.NewBuffer(post))
	if err != nil {
		return err
	}
	req.Header.Add("Connection", "close")
	req.Header.Add("User-Agent", variable.UserAgents[userAgentId])
	req.Header.Add("Host", XUST_HOST)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Referer", XUST_GET_COOKIE_URL+uid)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", "https://"+XUST_HOST)
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := (&http.Client{Transport: &http.Transport{DisableKeepAlives: true}}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if len(body) > 2 {
		return errors.New(string(body))
	}
	return nil
}

func XustSignInAfternoon() {
	logger := variable.Logger.With(zap.String("function", "XUST signin afternoon"))
	var results []model.XUSTUser
	err := variable.DB.Model(&model.XUSTUser{}).Where("time >= 2").Find(&results).Error
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, user := range results {
		logger := logger.With(zap.String("Sno", user.Sno))
		userAgentId := rand.Int() % len(variable.UserAgents)
		var info map[string]json.RawMessage
		var body []byte
		cookie, err := getXustCookie(user.Uid, userAgentId)
		if err != nil {
			goto NEXT_USER
		}
		info, err = getXustSigninInfo(user.Uid, user.Sno, cookie, userAgentId)
		if err != nil {
			goto NEXT_USER
		}
		body = getXustSignInRequestDate(info, XUST_AFTERNOON)
		for i := 3; i > 0; i-- {
			err = xustSignIn(user.Uid, cookie, userAgentId, body)
			if err != nil {
				continue
			}
			goto NEXT_USER
		}
	NEXT_USER:
		logger.Info(string(body))
		if err != nil {
			logger.Error(err.Error())
			err = SendEmail(
				mail.Address{Name: user.Sno, Address: user.Email}, "打卡失败",
				"今天的打卡任务<font color = 'red'>失败</font>了，请您手动打卡，谢谢！")
			if err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		logger.Info("SignIn Successfully")
	}
}

func XustSignInNight() {
	logger := variable.Logger.With(zap.String("function", "XUST signin night"))
	var results []model.XUSTUser
	err := variable.DB.Model(&model.XUSTUser{}).Find(&results).Error
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, user := range results {
		logger := logger.With(zap.String("Sno", user.Sno))
		userAgentId := rand.Int() % len(variable.UserAgents)
		var info map[string]json.RawMessage
		var body []byte
		cookie, err := getXustCookie(user.Uid, userAgentId)
		if err != nil {
			goto NEXT_USER
		}
		info, err = getXustSigninInfo(user.Uid, user.Sno, cookie, userAgentId)
		if err != nil {
			goto NEXT_USER
		}
		body = getXustSignInRequestDate(info, XUST_NIGHT)
		for i := 3; i > 0; i-- {
			err = xustSignIn(user.Uid, cookie, userAgentId, body)
			if err != nil {
				continue
			}
			goto NEXT_USER
		}
	NEXT_USER:
		logger.Info(string(body))
		if err != nil {
			logger.Error(err.Error())
			err = SendEmail(
				mail.Address{Name: user.Sno, Address: user.Email}, "打卡失败",
				"今天的打卡任务<font color = 'red'>失败</font>了，请您手动打卡，谢谢！")
			if err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		logger.Info("SignIn Successfully")
	}
}
