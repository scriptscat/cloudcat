package dto

import "github.com/scriptscat/cloudcat/internal/domain/sync/entity"

type SyncScript struct {
	Action     string             `json:"action"`
	Actiontime int64              `json:"actiontime"`
	UUID       string             `json:"uuid"`
	Msg        string             `json:"msg"`
	Script     *entity.SyncScript `json:"script"`
}

type SyncValue struct {
}
