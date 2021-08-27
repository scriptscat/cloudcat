package app

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/interface/http/apiv1"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"

	_ "github.com/scriptscat/cloudcat/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *config.Config) error {

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return err
	}
	kv, err := kvdb.NewKvDb(cfg.KvDB)

	r := gin.Default()

	if cfg.Mode != gin.ReleaseMode {
		url := ginSwagger.URL("/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	if err := apiv1.NewRouter(r, db, kv); err != nil {
		return err
	}

	return r.Run(cfg.Addr)
}
