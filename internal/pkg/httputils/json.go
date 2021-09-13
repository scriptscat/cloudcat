package httputils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/scriptscat/cloudcat/internal/controller/http/v1/dto"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
	pkgValidator "github.com/scriptscat/cloudcat/pkg/utils/validator"
	"github.com/sirupsen/logrus"
)

func Handle(ctx *gin.Context, f func() interface{}) {
	resp := f()
	if resp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok",
		})
		return
	}
	switch resp.(type) {
	case *errs.JsonRespondError:
		err := resp.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := resp.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := resp.(error)
		logrus.Errorf("%s - %s: %v", ctx.Request.RequestURI, ctx.ClientIP(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	case *dto.List:
		list := resp.(*dto.List)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "list": list.List, "total": list.Total,
		})
	default:
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "msg": "ok", "data": resp,
		})
	}
}

func HandleError(ctx *gin.Context, err error) {
	switch err.(type) {
	case *errs.JsonRespondError:
		err := err.(*errs.JsonRespondError)
		ctx.JSON(err.Status, err)
	case validator.ValidationErrors:
		err := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, "msg": pkgValidator.TransError(err),
		})
	case error:
		err := err.(error)
		logrus.Errorf("%s - %s: %v", ctx.Request.RequestURI, ctx.ClientIP(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, "msg": "系统错误",
		})
	}
}
