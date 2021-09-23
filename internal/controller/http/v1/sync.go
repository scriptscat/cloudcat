package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/internal/domain/sync/dto"
	"github.com/scriptscat/cloudcat/internal/domain/sync/service"
	"github.com/scriptscat/cloudcat/internal/pkg/httputils"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

type Sync struct {
	service.Sync
}

func NewSync() *Sync {
	return &Sync{}
}

// @Summary     同步
// @Description 推送脚本变更,需要先拉取获得版本号
// @ID          sync-push-script
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Param       version path int true "版本号"
// @Success     200 {object} dto.SyncScript
// @Failure     403
// @Router      /sync/{device}/script/push/{version} [put]
func (s *Sync) pushScript(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		sync := make([]*dto.SyncScript, 0)
		if err := c.BindJSON(&sync); err != nil {
			return err
		}
		version := utils.StringToInt64(c.Param("version"))
		ret, version, err := s.PushScript(uid, device, version, sync)
		if err != nil {
			return err
		}
		return gin.H{
			"version": version,
			"push":    ret,
		}
	})
}

// @Summary     同步
// @Description 拉取脚本变更
// @ID          sync-pull-script
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Param       version path int true "版本号"
// @Success     200 {object} dto.SyncScript
// @Failure     403
// @Router      /sync/{device}/script/pull/{version} [get]
func (s *Sync) pullScript(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		version := utils.StringToInt64(c.Param("version"))
		ret, version, err := s.PullScript(uid, device, version)
		if err != nil {
			return err
		}
		return gin.H{
			"version": version,
			"pull":    ret,
		}
	})
}

// @Summary     同步
// @Description 推送订阅变更,需要先拉取获得版本号
// @ID          sync-push-subscribe
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Param       version path int true "版本号"
// @Success     200 {object} dto.SyncSubscribe
// @Failure     403
// @Router      /sync/{device}/subscribe/push/{version} [put]
func (s *Sync) pushSubscribe(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		sync := make([]*dto.SyncSubscribe, 0)
		if err := c.BindJSON(&sync); err != nil {
			return err
		}
		version := utils.StringToInt64(c.Param("version"))
		ret, version, err := s.PushSubscribe(uid, device, version, sync)
		if err != nil {
			return err
		}
		return gin.H{
			"version": version,
			"push":    ret,
		}
	})
}

// @Summary     同步
// @Description 拉取脚本变更
// @ID          sync-pull-subscribe
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Param       version path int true "版本号"
// @Success     200 {object} dto.SyncSubscribe
// @Failure     403
// @Router      /sync/{device}/subscribe/pull/{version} [get]
func (s *Sync) pullSubscribe(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		version := utils.StringToInt64(c.Param("version"))
		ret, version, err := s.PullSubscribe(uid, device, version)
		if err != nil {
			return err
		}
		return gin.H{
			"version": version,
			"pull":    ret,
		}
	})
}

func (s *Sync) Register(r *gin.RouterGroup) {
	rg := r.Group("/sync/:device", userAuth())
	rgg := rg.Group("/script")
	rgg.PUT("/push/:version", s.pushScript)
	rgg.GET("/pull/:version", s.pullScript)

	rgg = rg.Group("/subscribe")
	rgg.PUT("/push/:version", s.pushSubscribe)
	rgg.GET("/pull/:version", s.pullSubscribe)

}
