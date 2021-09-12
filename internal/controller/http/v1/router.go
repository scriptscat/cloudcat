package v1

import (
	"github.com/gin-gonic/gin"
	repository2 "github.com/scriptscat/cloudcat/internal/domain/safe/repository"
	service3 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	service2 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
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
	systemConfig := config.NewSystemConfig(kv)
	userSvc := service.NewUser(systemConfig, kv, repository.NewBbsOAuth(db.DB), repository.NewWechatOAuth(db.DB), repository.NewUser(db.DB), repository.NewVerifyCode(kv))
	senderSvc := service2.NewSender(systemConfig)
	safeSvc := service3.NewSafe(repository2.NewSafe(kv))

	system := NewSystem(kv)

	user := NewUser(cfg.Jwt.Token, userSvc, safeSvc, senderSvc)

	register(v1, system, user)

	return nil
}
