package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/scriptscat/cloudcat/internal/controller/http/v1"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/migrations"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"

	_ "github.com/scriptscat/cloudcat/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config) error {

	db, err := database.NewDatabase(cfg.Database, cfg.Mode == "debug")
	if err != nil {
		return err
	}

	kv, err := kvdb.NewKvDb(cfg.KvDB)
	if err != nil {
		return err
	}

	if err := migrations.RunMigrations(db); err != nil {
		return err
	}

	binding.Validator = pkgValidator.NewValidator()

	r := gin.Default()

	if cfg.Mode != gin.ReleaseMode {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	if err := v1.NewRouter(r, cfg, db, kv); err != nil {
		return err
	}

	return r.Run(cfg.Addr)
}
