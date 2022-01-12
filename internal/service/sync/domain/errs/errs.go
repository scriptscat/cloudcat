package errs

import (
	"net/http"

	"github.com/scriptscat/cloudcat/pkg/errs"
)

var (
	ErrSyncVersionError = errs.NewBadRequestError(2001, "当前版本不是最新,请先同步数据!")
	ErrSyncIsNil        = errs.NewBadRequestError(2002, "同步数据是空的")

	ErrDeviceNotFound = errs.NewError(http.StatusNotFound, 2003, "设备不存在")
)
