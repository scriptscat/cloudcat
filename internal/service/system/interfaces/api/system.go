package api

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/infrastructure/kvdb"
	"github.com/scriptscat/cloudcat/internal/service/system/application"
	"github.com/scriptscat/cloudcat/internal/service/system/infrastructure/persistence"
	"github.com/scriptscat/cloudcat/pkg/httputils"
)

type System struct {
	application.System
}

func NewSystem(kv kvdb.KvDb) *System {
	return &System{
		System: application.NewSystem(persistence.NewRepo(kv)),
	}
}

// @Summary      系统
// @Description  查询脚本猫版本信息
// @ID           version
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  repository.ScriptCatInfo
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /system/version [get]
func (s *System) version(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		ret, err := s.ScriptCatInfo()
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary      系统环境
// @Description  获取系统环境变量
// @ID           env
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  repository.ScriptCatInfo
// @Failure      400  {object}  errs.JsonRespondError
// @Router       /system/version [get]
func (s *System) env(ctx *gin.Context) {
	httputils.Handle(ctx, func() interface{} {
		ret, err := s.ScriptCatInfo()
		if err != nil {
			return err
		}
		return ret
	})
}

func (s *System) Register(r *gin.RouterGroup) {
	v1 := r.Group("/system")
	v1.GET("/version", s.version)
	v1.GET("/env", s.env)
}
