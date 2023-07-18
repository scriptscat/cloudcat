package producer

import (
	"context"
	"encoding/json"

	"github.com/scriptscat/cloudcat/internal/model/entity/script_entity"

	"github.com/codfrm/cago/pkg/broker"
	broker2 "github.com/codfrm/cago/pkg/broker/broker"
)

// 脚本相关消息生产者

type ScriptUpdateMsg struct {
	Script *script_entity.Script
}

func PublishScriptUpdate(ctx context.Context, script *script_entity.Script) error {
	// code过大, 不传递
	script.Code = ""
	body, err := json.Marshal(&ScriptUpdateMsg{
		Script: script,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptUpdateTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptUpdateMsg(msg *broker2.Message) (*ScriptUpdateMsg, error) {
	ret := &ScriptUpdateMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeScriptUpdate(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptUpdateTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptUpdateMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m.Script)
	}, opts...)
	return err
}

type ScriptDeleteMsg struct {
	Script *script_entity.Script
}

func PublishScriptDelete(ctx context.Context, script *script_entity.Script) error {
	// code过大, 不传递
	script.Code = ""
	body, err := json.Marshal(&ScriptDeleteMsg{
		Script: script,
	})
	if err != nil {
		return err
	}
	return broker.Default().Publish(ctx, ScriptDeleteTopic, &broker2.Message{
		Body: body,
	})
}

func ParseScriptDeleteMsg(msg *broker2.Message) (*ScriptDeleteMsg, error) {
	ret := &ScriptDeleteMsg{}
	if err := json.Unmarshal(msg.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func SubscribeScriptDelete(ctx context.Context, fn func(ctx context.Context, script *script_entity.Script) error, opts ...broker2.SubscribeOption) error {
	_, err := broker.Default().Subscribe(ctx, ScriptDeleteTopic, func(ctx context.Context, ev broker2.Event) error {
		m, err := ParseScriptDeleteMsg(ev.Message())
		if err != nil {
			return err
		}
		return fn(ctx, m.Script)
	})
	return err
}
