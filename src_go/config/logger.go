package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"XustAutoSignIn/variable"
)

type log struct {
	Enabled bool   `yaml:"enabled"`                                                                   // 开启日志
	Level   string `yaml:"level,omitempty" validate:"required_if=Enabled true,oneof=info warn error"` // 日志等级
	Path    string `yaml:"path,omitempty" validate:"required_if=Enabled true,file"`                   // 日志路径
}

func (this *log) init() {
	if this.Enabled == false {
		return
	}
	var err error
	var level zapcore.Level
	switch this.Level {
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	if err != nil {
		variable.Logger.With(zap.String("feature", "log level")).Errorf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	f, err := os.OpenFile(this.Path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		variable.Logger.With(zap.String("feature", "create log file")).Fatalf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	_ = f.Close()
	variable.Logger = variable.NewFileLogger(this.Path, level)
}
