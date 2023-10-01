package subscribe

import (
	"context"

	"github.com/codfrm/cago/pkg/broker/broker"
	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"
	"github.com/scriptscat/cloudcat/internal/repository/resource_repo"
	"github.com/scriptscat/cloudcat/internal/task/producer"
)

type Resource struct {
}

func (v *Resource) Subscribe(ctx context.Context) error {
	if err := producer.SubscribeScriptUpdate(ctx, v.scriptUpdate, broker.Group("resource")); err != nil {
		return err
	}
	if err := producer.SubscribeScriptDelete(ctx, v.scriptDelete, broker.Group("resource")); err != nil {
		return err
	}
	return nil
}

// 消费脚本创建消息,根据meta信息进行分类
func (v *Resource) scriptUpdate(ctx context.Context, script *script_entity.Script) error {
	// 删除相关resource
	return v.scriptDelete(ctx, script)
}

func (v *Resource) scriptDelete(ctx context.Context, script *script_entity.Script) error {
	// 删除相关resource
	for _, v := range script.Metadata["require"] {
		if err := resource_repo.Resource().Delete(ctx, v); err != nil {
			return err
		}
	}
	return nil
}
