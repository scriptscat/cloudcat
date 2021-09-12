package errs

import (
	"fmt"
	"time"

	"github.com/scriptscat/cloudcat/internal/domain/safe/dto"
	"github.com/scriptscat/cloudcat/internal/pkg/errs"
)

func NewOperationTimeToShort(rule *dto.SafeRule) error {
	return errs.NewBadRequestError(4001, fmt.Sprintf("两次操作时间过断,请%d秒后重试", rule.Interval))
}

func NewOperationMax(rule *dto.SafeRule) error {
	return errs.NewBadRequestError(4002, fmt.Sprintf("%s,%d秒内重试了%d次,请%d秒后重试", rule.Description, rule.Period/time.Second, rule.Period, rule.Period/time.Second))
}

func NewOperationLimit(rule *dto.SafeRule) error {
	return errs.NewBadRequestError(4003, fmt.Sprintf("%s,已达上限,请%d秒后再重试", rule.Description, rule.Period/time.Second))
}
