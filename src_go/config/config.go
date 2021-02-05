package config

import (
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"XustAutoSignIn/flags"
	"XustAutoSignIn/variable"
)

type Config struct {
	Port           int       `yaml:"port,omitempty"`           // 服务器端口号，默认8080
	XustAutoSignIn bool      `yaml:"xustAutoSignIn,omitempty"` // 西科大自动打卡开关，默认false
	Database       *database `yaml:"database"`                 // 数据库配置
	Email          *email    `yaml:"email"`                    // 邮箱服务配置
	Log            *log      `yaml:"log,omitempty"`            // 日志
}

func GetConfig() *Config {
	configOnce.Do(func() {
		flag := flags.GetFlags()
		logger := variable.Logger.With(
			zap.String("module", "Config"),
		)
		configBytes, err := ioutil.ReadFile(flag.GetConfigFile())
		if err != nil {
			logger.With(zap.String("feature", "read config")).Fatalf(err.Error())
			os.Exit(variable.CONFIG_ERROR)
		}
		if err = yaml.Unmarshal(configBytes, &cfg); err != nil {
			logger.With(zap.String("feature", "unmarshal")).Fatalf(err.Error())
			os.Exit(variable.CONFIG_ERROR)
		}
		cfg.init(logger)
		flag.CheckTestFlag()
		logger.Infof("init config")
	})
	return cfg
}

func (this *Config) init(logger *zap.SugaredLogger) {
	if this.Port == 0 {
		this.Port = 8080
	}
	var err error
	// init database
	if err = variable.Validate.Struct(this.Database); err != nil {
		logger.With(zap.String("feature", "check database config")).Errorf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	this.Database.init()
	// init smtp client
	if err = variable.Validate.Struct(this.Email); err != nil {
		logger.With(zap.String("feature", "check email config")).Errorf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	this.Email.init()
	// init logger
	if cfg.Log == nil || cfg.Log.Enabled == false {
		return
	}
	if err = variable.Validate.Struct(this.Log); err != nil {
		logger.With(zap.String("feature", "check logger config")).Errorf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	this.Log.init()
}

func (this *Config) GetEmailClient() *smtp.Client {
	return this.Email.getClient()
}

func (this *Config) GetFromEmail() mail.Address {
	return mail.Address{
		Name:    this.Email.User,
		Address: this.Email.From,
	}
}

func (this *Config) GetServerPort() int {
	return this.Port
}

var (
	cfg        *Config
	configOnce sync.Once
)
