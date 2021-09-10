package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type Register interface {
	Register(r *gin.RouterGroup)
}

func register(r *gin.RouterGroup, register ...Register) {
	for _, v := range register {
		v.Register(r)
	}
}

// Swagger spec:
// @title       云猫api文档
// @version     1.0
// @BasePath    /api/v1

func NewRouter(r *gin.Engine, cfg *config.Config, db *database.Database, kv kvdb.KvDb) error {

	v1 := r.Group("/api/v1")

	userSvc := service.NewUser(config.NewSystemConfig(kv))

	system := NewSystem(kv)

	user := NewUser(cfg.Jwt.Token, userSvc)

	register(v1, system, user)

	return nil
}
