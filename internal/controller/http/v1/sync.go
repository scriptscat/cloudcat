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
// @Description 同步脚本
// @ID          sync-script
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Success     200 {object} []*dto.SyncScript
// @Failure     403
// @Router      /sync/{device}/script/sync [put]
func (s *Sync) syncScript(c *gin.Context) {
	httputils.Handle(c, func() interface{} {
		uid, _ := userId(c)
		device := utils.StringToInt64(c.Param("device"))
		sync := make([]*dto.SyncScript, 0)
		if err := c.BindJSON(&sync); err != nil {
			return err
		}
		ret, err := s.SyncScript(uid, device, sync)
		if err != nil {
			return err
		}
		return ret
	})
}

// @Summary     同步
// @Description 同步脚本
// @ID          sync-script
// @Tags  	    sync
// @Accept      json
// @Security    BearerAuth
// @Param       device path int true "设备id"
// @Success     200 {object} []*dto.SyncScript
// @Failure     403
// @Router      /sync/{device}/script/sync [put]
func (s *Sync) fullScript(c *gin.Context) {

}

func (s *Sync) Register(r *gin.RouterGroup) {
	rg := r.Group("/sync/:device", tokenAuth())
	rgg := rg.Group("/script")
	rgg.PUT("/sync", s.syncScript)
	rgg.PUT("/full", s.syncScript)

}
