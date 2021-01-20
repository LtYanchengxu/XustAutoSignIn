package config

import (
	"crypto/tls"
	"net/smtp"
	"strconv"

	"go.uber.org/zap"

	"XustAutoSignIn/variable"
)

type email struct {
	Host     string `yaml:"host" validate:"required"`                  // 服务器地址
	Port     int    `yaml:"port,omitempty" validate:"min=0,max=65535"` // SMTP 端口，默认25
	Account  string `yaml:"account" validate:"required"`               // 账户
	Password string `yaml:"password"`                                  // 密码
	User     string `yaml:"username,omitempty"`                        // 用户名
	From     string `yaml:"from" validate:"email|eqfield=Account"`     // From Email
	Starttls bool   `yaml:"starttls,omitempty"`                        // 开启 starttls，默认 false
	Ssl      bool   `yaml:"ssl,omitempty"`                             // 开启 TLS，默认 false
	// Smime          *smime  `yaml:"smime"`          // S/MIME 配置，默认为空，不开启
}

type smime struct {
	Enabled bool   `yaml:"enabled,omitempty"` // 开启 S/MIME，默认 false
	Cert    string `yaml:"cert,omitempty"`    // S/MIME 证书文件，默认 ./.cert_file
	Key     string `yaml:"key,omitempty"`     // S/MIME 密钥文件，默认 ./.key_file
}

func (this *email) init() {
	if this.Port == 0 {
		this.Port = 25
	}
}

func (this *email) getClient() (client *smtp.Client) {
	server := this.Host + ":" + strconv.Itoa(this.Port)
	logger := variable.Logger.With(zap.String("module", "Config"), zap.String("server", server))
	var err error
	tlsConfig := &tls.Config{
		ServerName:         this.Host,
		InsecureSkipVerify: true,
	}
	// SMTP Client
	if this.Ssl == true {
		var conn *tls.Conn
		conn, err = tls.Dial("tcp", server, tlsConfig)
		if err != nil {
			logger.With(zap.String("feature", "tls connection")).Error(err.Error())
			return nil
		}
		client, err = smtp.NewClient(conn, this.Host)
	} else {
		client, err = smtp.Dial(server)
	}
	if err != nil {
		logger.With(zap.String("feature", "smtp client")).Error(err.Error())
		return nil
	}
	// StartTLS
	if this.Starttls == true {
		if err = client.StartTLS(tlsConfig); err != nil {
			logger.With(zap.String("feature", "starttls")).Error(err.Error())
			client.Close()
			return nil
		}
	}
	// Config
	if err = client.Auth(smtp.PlainAuth("", this.Account, this.Password, this.Host)); err != nil {
		logger.With(zap.String("feature", "email auth")).Error(err.Error())
		client.Close()
		return nil
	}
	if err = client.Mail(this.From); err != nil {
		logger.With(zap.String("feature", "email from")).Error(err.Error())
		client.Close()
		return nil
	}
	return client
}
