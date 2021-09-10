package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/domain/system/repository"
	"github.com/scriptscat/cloudcat/internal/domain/system/service"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/kvdb"
)

type System struct {
	service.System
}

func NewSystem(kv kvdb.KvDb) *System {
	return &System{
		System: service.NewSystem(repository.NewRepo(kv)),
	}
}

// @Summary     系统
// @Description 查询脚本猫版本信息
// @ID          system
// @Tags  	    system
// @Accept      json
// @Produce     json
// @Success     200 {object} repository.ScriptCatInfo
// @Failure     400 {object} errs.JsonRespondError
// @Router      /system/version [get]
func (s *System) version(ctx *gin.Context) {
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
}
