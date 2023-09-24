package subscribe

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
	"github.com/scriptscat/cloudcat/internal/task/producer"
)

type Value struct {
}

func (v *Value) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptUpdate(ctx, v.scriptUpdate, broker.Group("value")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptDelete(ctx, v.scriptDelete, broker.Group("value")); err != nil {
		return err
	}
	return nil
}

// 消费脚本创建消息,根据meta信息进行分类
func (v *Value) scriptUpdate(ctx context.Context, script *script_entity.Script) error {
	return nil
}

func (v *Value) scriptDelete(ctx context.Context, script *script_entity.Script) error {
	// 查询还有哪些脚本使用了这个storage
	list, err := script_repo.Script().FindByStorage(ctx, script.StorageName())
	if err != nil {
		return err
	}
	if len(list) == 0 {
		// 删除storageName下的value
		if err := value_repo.Value().DeleteByStorage(ctx, script.StorageName()); err != nil {
			return err
		}
		// 删除storageCookie下的value
		if err := cookie_repo.Cookie().DeleteByStorage(ctx, script.StorageName()); err != nil {
			return err
		}
	}
	return nil
}
