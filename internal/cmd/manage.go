package cmd

import (
	"github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/entity"
	service2 "github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/utils"
	"github.com/spf13/cobra"
)

type manageCmd struct {
	config string
	user   string
	passwd string
	email  string
	host   string
	tls    bool

	platform     string
	clientId     string
	clientSecret string
	token        string
	aes          string
	appId        string
	appSecret    string
}

func NewManageCmd() *manageCmd {
	return &manageCmd{}
}

func (m *manageCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "manage [flag]",
		Short: "管理系统配置等信息",
	}
	ret.Flags().StringVarP(&m.config, "config", "c", "config.yaml", "配置文件")
	admin := &cobra.Command{
		Use:   "admin [flag]",
		Short: "修改admin用户信息",
		RunE:  m.admin,
	}
	admin.Flags().StringVarP(&m.user, "user", "u", "", "管理员账号")
	admin.Flags().StringVarP(&m.passwd, "passwd", "p", "", "管理员密码")
	admin.Flags().StringVarP(&m.email, "email", "e", "", "管理员邮箱")

	sender := &cobra.Command{
		Use:   "sender [flag]",
		Short: "配置验证码发送",
		RunE:  m.sender,
	}

	sender.Flags().StringVarP(&m.host, "host", "", "", "邮箱服务器地址(eg:smtp.scriptcat.org:587)")
	sender.Flags().StringVarP(&m.email, "email", "e", "", "邮箱发件账号")
	sender.Flags().StringVarP(&m.passwd, "passwd", "p", "", "邮箱发件密码")
	sender.Flags().BoolVarP(&m.tls, "tls", "t", false, "是否使用tls链接")

	oauth := &cobra.Command{
		Use:   "oauth [flag]",
		Short: "配置三方登录信息",
		RunE:  m.oauth,
	}

	oauth.Flags().StringVarP(&m.platform, "platform", "", "", "设置平台(wechat,bbs)")
	oauth.Flags().StringVarP(&m.clientId, "client-id", "", "", "论坛登录clientId")
	oauth.Flags().StringVarP(&m.clientSecret, "client-secret", "", "", "论坛登录clientSecret")
	oauth.Flags().StringVarP(&m.appId, "app-id", "", "", "微信AppId")
	oauth.Flags().StringVarP(&m.appSecret, "app-secret", "", "", "微信AppSecret")
	oauth.Flags().StringVarP(&m.token, "token", "", "", "微信Token")
	oauth.Flags().StringVarP(&m.aes, "aes", "", "", "微信aes密钥")

	ret.AddCommand(admin, sender, oauth)
	return []*cobra.Command{ret}
}

func (m *manageCmd) admin(cmd *cobra.Command, args []string) error {
	cfg, err := config.Init(m.config)
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(cfg.Database, cfg.Mode == "debug")
	if err != nil {
		return err
	}
	user := &entity.User{}
	if err := db.Where("id=1").First(&user).Error; err != nil {
		return err
	}
	user.Nickname = m.getValue(m.user, user.Nickname)
	if m.passwd != "" {
		if err := user.SetPassword(m.passwd); err != nil {
			return err
		}
	}
	user.Email = m.getValue(m.email, user.Email)
	return db.Save(user).Error
}

func (m *manageCmd) sender(cmd *cobra.Command, args []string) error {
	config, err := m.getConfig()
	if err != nil {
		return err
	}

	if m.tls {
		err = config.SetConfig(service.SENDER_EMAIL_TLS, "1")
	} else {
		err = config.SetConfig(service.SENDER_EMAIL_TLS, "0")
	}

	return utils.Errs(
		err,
		config.SetConfig(service.SENDER_EMAIL_HOST, m.host),
		config.SetConfig(service.SENDER_EMAIL_USER, m.email),
		config.SetConfig(service.SENDER_EMAIL_PASSWD, m.passwd),
	)
}

func (m *manageCmd) oauth(cmd *cobra.Command, args []string) error {
	config, err := m.getConfig()
	if err != nil {
		return err
	}
	errs := make([]error, 0)
	switch m.platform {
	case "wechat":
		errs = append(errs, config.SetConfig(service2.OAuthConfigWechatAppId, m.appId))
		errs = append(errs, config.SetConfig(service2.OAuthConfigWechatAppSecret, m.appSecret))
		errs = append(errs, config.SetConfig(service2.OAuthConfigWechatToken, m.token))
		errs = append(errs, config.SetConfig(service2.OAuthConfigWechatEncodingaeskey, m.aes))
	case "bbs":
		errs = append(errs, config.SetConfig(service2.OAuthConfigBbsClientId, m.clientId))
		errs = append(errs, config.SetConfig(service2.OAuthConfigBbsClientSecret, m.clientSecret))
	}
	return utils.Errs(errs...)
}

func (m *manageCmd) getConfig() (config.SystemConfig, error) {
	cfg, err := config.Init(m.config)
	if err != nil {
		return nil, err
	}

	kv, err := kvdb.NewKvDb(cfg.KvDB)
	if err != nil {
		return nil, err
	}

	config := config.NewSystemConfig(kv)
	return config, nil
}

func (m *manageCmd) getValue(value, defVal string) string {
	if value == "" {
		return defVal
	}
	return value
}
