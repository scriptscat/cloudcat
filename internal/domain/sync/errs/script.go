package errs

import "github.com/scriptscat/cloudcat/internal/pkg/errs"

var (
	ErrSyncVersionError = errs.NewBadRequestError(2001, "当前版本不是最新,请先同步数据!")
	ErrSyncIsNil        = errs.NewBadRequestError(2002, "同步数据是空的")
)
