package server

import (
	"context"
	"fmt"
	"github.com/codfrm/cago/pkg/broker"
	"github.com/scriptscat/cloudcat/internal/api/auth"
	"github.com/scriptscat/cloudcat/internal/repository/token_repo"
	"github.com/scriptscat/cloudcat/internal/service/auth_svc"
	"github.com/scriptscat/cloudcat/internal/task/consumer"
	"github.com/scriptscat/cloudcat/migrations"
	"github.com/scriptscat/cloudcat/pkg/cloudcat_api"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"

	"github.com/scriptscat/cloudcat/pkg/bbolt"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/api"
	"github.com/spf13/cobra"
)

type Server struct {
	config *string
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Command(config *string) []*cobra.Command {
	server := &cobra.Command{
		Use:   "server",
		Short: "启动服务",
		RunE:  s.Server,
	}

	init := &cobra.Command{
		Use:   "init",
		Short: "初始化服务",
		RunE:  s.Init,
	}

	s.config = config
	return []*cobra.Command{server, init}
}

func (s *Server) Server(cmd *cobra.Command, args []string) error {
	cfg, err := configs.NewConfig("cloudcat", configs.WithConfigFile(*s.config))
	if err != nil {
		log.Fatalf("new config err: %v", err)
	}
	err = cago.New(cmd.Context(), cfg).DisableLogger().
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(bbolt.Bolt()).
		Registry(cago.FuncComponent(consumer.Consumer)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(bbolt.Default())
		})).
		RegistryCancel(mux.HTTP(api.Router)).
		Start()
	if err != nil {
		return err
	}
	return nil
}

var configTemplate = `version: 1.0.0
source: file
broker:
    type: "event_bus"
debug: false
env: dev
http:
    address:
        - :8644
logger:
    level: "info"
    logfile:
        enable: true
        errorfilename: "{{.configDir}}/cloudcat.error.log"
        filename: "{{.configDir}}/cloudcat.log"
db:
  path: "{{.configDir}}/data.db"
`

func (s *Server) Init(cmd *cobra.Command, args []string) error {
	// 判断文件是否存在
	_, err := configs.NewConfig("cloudcat", configs.WithConfigFile(*s.config))
	if err == nil {
		return fmt.Errorf("config file %s is exist", *s.config)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("determine whether the file exists: %w", err)
	}
	// 写配置文件
	configDir := path.Dir(*s.config)
	if err := os.MkdirAll(configDir, 0744); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	// 模板渲染
	tpl := strings.ReplaceAll(configTemplate, "{{.configDir}}", configDir)
	if err := os.WriteFile(*s.config, []byte(tpl), 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	cfg, err := configs.NewConfig("cloudcat", configs.WithConfigFile(*s.config))
	if err != nil {
		return fmt.Errorf("new config err: %w", err)
	}
	err = cago.New(cmd.Context(), cfg).DisableLogger().
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(broker.Broker)).
		Registry(bbolt.Bolt()).
		Registry(cago.FuncComponent(consumer.Consumer)).
		Registry(cago.FuncComponent(func(ctx context.Context, cfg *configs.Config) error {
			return migrations.RunMigrations(bbolt.Default())
		})).
		RegistryCancel(cago.FuncComponentCancel(func(ctx context.Context, cancel context.CancelFunc, cfg *configs.Config) error {
			defer cancel()
			token_repo.RegisterToken(token_repo.NewToken())
			// 写client配置文件
			config := &cloudcat_api.Config{
				ApiVersion: "v1",
				Server: &cloudcat_api.ConfigServer{
					BaseURL: "http://127.0.0.1:8644",
				},
			}
			token, err := auth_svc.Token().TokenCreate(ctx, &auth.TokenCreateRequest{
				TokenID: "default",
			})
			if err != nil {
				return fmt.Errorf("create token err: %w", err)
			}
			config.User = &cloudcat_api.ConfigUser{
				Name:              token.Token.ID,
				Token:             token.Token.Token,
				DataEncryptionKey: token.Token.DataEncryptionKey,
			}
			data, err := yaml.Marshal(config)
			if err != nil {
				return fmt.Errorf("marshal config err: %w", err)
			}
			if err := os.WriteFile(path.Join(configDir, "cloudcat.yaml"), data, 0644); err != nil {
				return fmt.Errorf("write client config file: %w", err)
			}
			return nil
		})).
		Start()
	if err != nil {
		return err
	}
	return nil

	return nil
}
