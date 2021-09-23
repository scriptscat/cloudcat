package dto

import "github.com/scriptscat/cloudcat/internal/domain/sync/entity"

type SyncScript struct {
	Action     string             `json:"action"`
	Actiontime int64              `json:"actiontime"`
	UUID       string             `json:"uuid"`
	Msg        string             `json:"msg,omitempty"`
	Script     *entity.SyncScript `json:"script,omitempty"`
}

type SyncSubscribe struct {
	Action     string                `json:"action"`
	Actiontime int64                 `json:"actiontime"`
	URL        string                `json:"url"`
	Msg        string                `json:"msg,omitempty"`
	Subscribe  *entity.SyncSubscribe `json:"subscribe,omitempty"`
}

type SyncValue struct {
}
