package v1

import (
	"os"

	"github.com/gin-gonic/gin"
	repository2 "github.com/scriptscat/cloudcat/internal/domain/safe/repository"
	service3 "github.com/scriptscat/cloudcat/internal/domain/safe/service"
	repository3 "github.com/scriptscat/cloudcat/internal/domain/sync/repository"
	service4 "github.com/scriptscat/cloudcat/internal/domain/sync/service"
	service2 "github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/domain/user/repository"
	"github.com/scriptscat/cloudcat/internal/domain/user/service"
	"github.com/scriptscat/cloudcat/internal/pkg/config"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/cache"
	"github.com/scriptscat/cloudcat/pkg/database"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
	"github.com/scriptscat/cloudcat/pkg/middleware/token"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

type Register interface {
	Register(r *gin.RouterGroup)
}

func register(r *gin.RouterGroup, register ...Register) {
	for _, v := range register {
		v.Register(r)
	}
}

var tokenAuth func(enforce bool) func(ctx *gin.Context)
var userAuth func() func(ctx *gin.Context)

// NewRouter 初始化路由
// Swagger spec:
// @title       云猫api文档
// @version     1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath    /api/v1
func NewRouter(r *gin.Engine, cfg *config.Config, db *database.Database, kv kvdb.KvDb, cache cache.Cache) error {
	disableAuth := os.Getenv("DISABLE_AUTH") == "true"
	tokenAuth = func(enforce bool) func(ctx *gin.Context) {
		if disableAuth {
			return token.Middleware(cache, enforce, token.WithExpired(TokenAuthMaxAge), token.WithDebug(gin.H{
				"uid":      "1",
				"username": "admin",
				"token":    utils.RandString(16, 1),
			}))
		}
		return token.Middleware(cache, enforce, token.WithExpired(TokenAuthMaxAge))
	}

	v1 := r.Group("/api/v1")
	systemConfig := config.NewSystemConfig(kv)
	userSvc := service.NewUser(systemConfig, kv, repository.NewUser(db.DB), repository.NewVerifyCode(kv))
	oauthSvc := service.NewOAuth(systemConfig, kv, db.DB, userSvc, repository.NewBbsOAuth(db.DB), repository.NewWechatOAuth(db.DB, kv))
	senderSvc := service2.NewSender(systemConfig)
	safeSvc := service3.NewSafe(repository2.NewSafe(kv))
	syncSvc := service4.NewSync(repository3.NewDevice(db.DB), repository3.NewScript(db.DB, kv), repository3.NewSubscribe(db.DB, kv))

	system := NewSystem(kv)

	auth := NewAuth(cache, userSvc, oauthSvc, safeSvc, senderSvc)
	user := NewUser(userSvc, oauthSvc, safeSvc, senderSvc)
	sync := NewSync(syncSvc)

	enforceAuth := tokenAuth(true)
	if disableAuth {
		userAuth = func() func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				enforceAuth(ctx)
			}
		}
	} else {
		userAuth = func() func(ctx *gin.Context) {
			return func(ctx *gin.Context) {
				enforceAuth(ctx)
				if !ctx.IsAborted() {
					uid, _ := userId(ctx)
					if _, err := auth.UserInfo(uid); err != nil {
						httputils.HandleError(ctx, err)
						ctx.Abort()
					}
				}
			}
		}
	}

	register(v1, system, user, auth, sync)

	return nil
}
