package application

import (
	"crypto/tls"
	"errors"
	"strconv"
	"strings"

	"github.com/scriptscat/cloudcat/internal/infrastructure/config"
	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendEmail(to, title, content, contextType string) error
}

const (
	SENDER_EMAIL_HOST   = "sender_email_host"
	SENDER_EMAIL_USER   = "sender_email_user"
	SENDER_EMAIL_PASSWD = "sender_email_passwd"
	SENDER_EMAIL_TLS    = "sender_email_tls"
)

type sender struct {
	config config.SystemConfig
}

func NewSender(config config.SystemConfig) Sender {
	return &sender{
		config: config,
	}
}

func (s *sender) SendEmail(to, title, content, contextType string) error {
	host, err := s.config.GetConfig(SENDER_EMAIL_HOST)
	if err != nil {
		return err
	}
	if host == "" {
		return errors.New("邮件发送器配置错误")
	}
	hosts := strings.SplitN(host, ":", 2)
	port, _ := strconv.Atoi(hosts[1])
	user, _ := s.config.GetConfig(SENDER_EMAIL_USER)
	passwd, _ := s.config.GetConfig(SENDER_EMAIL_PASSWD)
	tlsCfg, _ := s.config.GetConfig(SENDER_EMAIL_TLS)

	d := gomail.NewDialer(hosts[0], port, user, passwd)
	if tlsCfg == "1" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	mark, err := s.config.GetConfig(SYSTEM_SITE_NAME)
	if err != nil {
		return err
	}
	if mark == "" {
		mark = "脚本猫"
	}
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "["+mark+"]"+title)
	m.SetBody(contextType, content)

	return d.DialAndSend(m)
}
