package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/infrastructure/config"
	"github.com/scriptscat/cloudcat/internal/infrastructure/middleware/token"
	"github.com/scriptscat/cloudcat/internal/infrastructure/persistence"
	"github.com/scriptscat/cloudcat/internal/infrastructure/sender"
	"github.com/scriptscat/cloudcat/internal/pkg/database"
	"github.com/scriptscat/cloudcat/internal/pkg/kvdb"
	service3 "github.com/scriptscat/cloudcat/internal/service/safe/application"
	repository2 "github.com/scriptscat/cloudcat/internal/service/safe/domain/repository"
	service4 "github.com/scriptscat/cloudcat/internal/service/sync/application"
	repository3 "github.com/scriptscat/cloudcat/internal/service/sync/domain/repository"
	api2 "github.com/scriptscat/cloudcat/internal/service/sync/interfaces/api"
	"github.com/scriptscat/cloudcat/internal/service/system/interfaces/api"
	application2 "github.com/scriptscat/cloudcat/internal/service/user/application"
	repository5 "github.com/scriptscat/cloudcat/internal/service/user/domain/repository"
	api3 "github.com/scriptscat/cloudcat/internal/service/user/interfaces/api"
	"github.com/scriptscat/cloudcat/pkg/cache"
	"github.com/scriptscat/cloudcat/pkg/httputils"
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

// NewRouter 初始化路由
// Swagger spec:
// @title                       云猫api文档
// @version                     1.0
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @BasePath                    /api/v1
func NewRouter(r *gin.Engine, db *database.Database, kv kvdb.KvDb, cache cache.Cache, repo *persistence.Repositories) error {
	disableAuth := os.Getenv("DISABLE_AUTH") == "true"
	token.TokenAuth = func(enforce bool) func(ctx *gin.Context) {
		if disableAuth {
			return token.Middleware(cache, enforce, token.WithExpired(api3.TokenAuthMaxAge), token.WithDebug(gin.H{
				"uid":      "1",
				"username": "admin",
				"token":    utils.RandString(16, 1),
			}))
		}
		return token.Middleware(cache, enforce, token.WithExpired(api3.TokenAuthMaxAge))
	}

	v1 := r.Group("/api/v1")
	systemConfig := config.NewSystemConfig(kv)
	senderSvc := sender.NewSender(systemConfig)
	userSvc := application2.NewUser(systemConfig, kv, repo.User.User, repo.User.VerifyCode, senderSvc)
	oauthSvc := application2.NewOAuth(systemConfig, kv, db.DB, userSvc, repository5.NewBbsOAuth(db.DB), repository5.NewWechatOAuth(db.DB, kv))
	safeSvc := service3.NewSafe(repository2.NewSafe(kv))
	syncSvc := service4.NewSync(repository3.NewDevice(db.DB), repository3.NewScript(db.DB, kv), repository3.NewSubscribe(db.DB, kv))

	system := api.NewSystem(kv)

	auth := api3.NewAuth(cache, userSvc, oauthSvc, safeSvc, senderSvc)
	user := api3.NewUser(userSvc, oauthSvc, safeSvc, senderSvc)
	sync := api2.NewSync(syncSvc)

	if disableAuth {
		token.UserAuth = func(enforce bool) func(ctx *gin.Context) {
			auth := token.TokenAuth(enforce)
			return func(ctx *gin.Context) {
				auth(ctx)
			}
		}
	} else {
		token.UserAuth = func(enforce bool) func(ctx *gin.Context) {
			authHandler := token.TokenAuth(enforce)
			return func(ctx *gin.Context) {
				authHandler(ctx)
				if !ctx.IsAborted() {
					uid, _ := token.UserId(ctx)
					if uid != 0 {
						// NOTE:用户信息可以写入context
						if _, err := auth.UserInfo(uid); err != nil {
							httputils.HandleError(ctx, err)
							ctx.Abort()
						}
					}
				}
			}
		}
	}

	register(v1, system, user, auth, sync)

	return nil
}
