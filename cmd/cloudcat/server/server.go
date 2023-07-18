package server

import (
	"context"
	"log"

	"github.com/codfrm/cago/pkg/broker"
	"github.com/scriptscat/cloudcat/internal/task/consumer"
	"github.com/scriptscat/cloudcat/migrations"

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
	ret := &cobra.Command{
		Use:   "server",
		Short: "启动服务",
		RunE:  s.Server,
	}
	s.config = config
	return []*cobra.Command{ret}
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
