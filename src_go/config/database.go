package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"XustAutoSignIn/model"
	"XustAutoSignIn/variable"
)

type database struct {
	Type     string `yaml:"type" validate:"required,oneof=mysql pgsql sqlite"` // 数据库类型
	Host     string `yaml:"host,omitempty" validate:"ip|file"`                 // 数据库地址 IP，默认 127.0.0.1
	Port     int    `yaml:"port,omitempty" validate:"min=0,max=65535"`         // 数据库端口，默认 mysql(3306) pgsql(5432)
	User     string `yaml:"user,omitempty"`                                    // 数据库用户
	Password string `yaml:"password,omitempty"`                                // 数据库密码，默认为空
	Name     string `yaml:"name,omitempty"`                                    // 数据库名称
}

func (this *database) init() {
	if this.Host == "" {
		this.Host = "127.0.0.1"
	}
	if this.User == "" {
		this.User = "root"
	}
	if this.Name == "" {
		this.Name = "test"
	}
	logger := variable.Logger.With(
		zap.String("module", "Config"),
	)
	var err error
	switch this.Type {
	case "mysql":
		if this.Port == 0 {
			this.Port = 3306
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			this.User, this.Password, this.Host, this.Port, this.Name)
		variable.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "pgsql":
		if this.Port == 0 {
			this.Port = 5432
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			this.Host, this.User, this.Password, this.Name, this.Port)
		variable.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		variable.DB, err = gorm.Open(sqlite.Open(this.Host), &gorm.Config{})
	}
	if err != nil {
		logger.With(zap.String("feature", "database connection")).Fatalf(err.Error())
		os.Exit(variable.CONFIG_ERROR)
	}
	this.migration()
}

func (this *database) migration() {
	variable.DB.AutoMigrate(&model.XUSTUser{})
}
