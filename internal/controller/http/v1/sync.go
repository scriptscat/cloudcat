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

func NewSync(sync service.Sync) *Sync {
	return &Sync{
		Sync: sync,
	}
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
// @Description 拉取订阅变更
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

// @Summary     同步
// @Description 获取设备列表
// @ID          sync-device-list
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Success     200 {object} entity.SyncDevice
// @Failure     403
// @Router      /sync/device [get]
func (s *Sync) device(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		list, err := s.Sync.DeviceList(uid)
		if err != nil {
			return err
		}
		return list
	})
}

// @Summary     同步
// @Description 推送设置变更
// @ID          sync-push-setting
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Param       setting formData string true "设置json字符串"
// @Param       settingtime formData string true "设置更新时间"
// @Success     200
// @Failure     403
// @Router      /sync/{device}/setting/push [put]
func (s *Sync) pushSetting(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		setting := c.PostForm("setting")
		settingtime := utils.StringToInt64(c.PostForm("settingtime"))
		err := s.PushSetting(uid, device, setting, settingtime)
		if err != nil {
			return err
		}
		return nil
	})
}

// @Summary     同步
// @Description 拉取设置变更
// @ID          sync-pull-setting
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Success     200
// @Failure     403
// @Router      /sync/{device}/setting/pull [get]
func (s *Sync) pullSetting(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		ret, time, err := s.PullSetting(uid, device)
		if err != nil {
			return err
		}
		return gin.H{
			"setting":     ret,
			"settingtime": time,
		}
	})
}

func (s *Sync) Register(r *gin.RouterGroup) {
	rg := r.Group("/sync", userAuth(true))
	rg.GET("/device", s.device)

	rg = rg.Group("/:device")

	rg.PUT("/setting/push", s.pushSetting)
	rg.GET("/setting/pull", s.pullSetting)

	rgg := rg.Group("/script")
	rgg.PUT("/push/:version", s.pushScript)
	rgg.GET("/pull/:version", s.pullScript)

	rgg = rg.Group("/subscribe")
	rgg.PUT("/push/:version", s.pushSubscribe)
	rgg.GET("/pull/:version", s.pullSubscribe)

}
