package v1

import (
	"github.com/gin-gonic/gin"
	repository2 "github.com/scriptscat/cloudcat/internal/domain/safe/repository"
	service3 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	service2 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/middleware/jwt"
)

type Register interface {
	Register(r *gin.RouterGroup)
}

func register(r *gin.RouterGroup, register ...Register) {
	for _, v := range register {
		v.Register(r)
	}
}

var jwtAuth func(enforce bool) func(ctx *gin.Context)
var tokenAuth func() func(ctx *gin.Context)

// NewRouter
// Swagger spec:
// @title       云猫api文档
// @version     1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath    /api/v1
func NewRouter(r *gin.Engine, cfg *config.Config, db *database.Database, kv kvdb.KvDb) error {
	jwtAuth = func(enforce bool) func(ctx *gin.Context) {
		return jwt.Jwt([]byte(cfg.Jwt.Token), enforce, jwt.WithExpired(JwtAuthMaxAge))
	}

	v1 := r.Group("/api/v1")
	systemConfig := config.NewSystemConfig(kv)
	userSvc := service.NewUser(systemConfig, kv, repository.NewUser(db.DB), repository.NewVerifyCode(kv))
	oauthSvc := service.NewOAuth(systemConfig, kv, db.DB, userSvc, repository.NewBbsOAuth(db.DB), repository.NewWechatOAuth(db.DB, kv))
	senderSvc := service2.NewSender(systemConfig)
	safeSvc := service3.NewSafe(repository2.NewSafe(kv))

	system := NewSystem(kv)

	auth := NewAuth(cfg.Jwt.Token, userSvc, oauthSvc, safeSvc, senderSvc)
	user := NewUser(cfg.Jwt.Token, userSvc, oauthSvc, safeSvc, senderSvc)

	enforceJwt := jwt.Jwt([]byte(cfg.Jwt.Token), true, jwt.WithExpired(JwtAuthMaxAge))
	tokenAuth = func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			enforceJwt(ctx)
			if !ctx.IsAborted() {
				uid, _ := userId(ctx)
				if _, err := auth.UserInfo(uid); err != nil {
					httputils.HandleError(ctx, err)
					ctx.Abort()
				}
			}
		}
	}

	register(v1, system, user, auth)

	return nil
}
