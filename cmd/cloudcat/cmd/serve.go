package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/scriptscat/cloudcat/internal/infrastructure/config"
	"github.com/scriptscat/cloudcat/internal/infrastructure/database"
	"github.com/scriptscat/cloudcat/internal/infrastructure/kvdb"
	v1 "github.com/scriptscat/cloudcat/internal/interfaces/api"
	"github.com/scriptscat/cloudcat/migrations"
	cache2 "github.com/scriptscat/cloudcat/pkg/cache"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
	"github.com/spf13/cobra"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type serveCmd struct {
	config string
}

func NewServeCmd() *serveCmd {
	return &serveCmd{}
}

func (s *serveCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "serve [flag]",
		Short: "运行脚本猫服务",
		RunE:  s.serve,
	}
	ret.Flags().StringVarP(&s.config, "config", "c", "config.yaml", "配置文件")

	return []*cobra.Command{ret}
}

func (s *serveCmd) serve(cmd *cobra.Command, args []string) error {
	cfg, err := config.Init(s.config)
	if err != nil {
		return err
	}

	if err := s.run(cfg); err != nil {
		return fmt.Errorf("serve start err: %v", err)
	}
	return nil
}

func (s *serveCmd) run(cfg *config.Config) error {
	db, err := database.NewDatabase(cfg.Database, cfg.Mode == "debug")
	if err != nil {
		return err
	}
	kv, err := kvdb.NewKvDb(cfg.KvDB)
	if err != nil {
		return err
	}
	cache, err := cache2.NewCache(cfg.Cache)
	if err != nil {
		return err
	}
	if err := migrations.RunMigrations(db); err != nil {
		return err
	}

	binding.Validator = pkgValidator.NewValidator()

	gin.SetMode(cfg.Mode)
	r := gin.Default()

	if cfg.Mode != gin.ReleaseMode {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	if err := v1.NewRouter(r, cfg, db, kv, cache); err != nil {
		return err
	}

	return r.Run(cfg.Addr)
}
