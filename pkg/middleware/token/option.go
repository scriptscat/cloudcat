package token

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Token struct {
	Info       gin.H  `json:"info"`
	Token      string `json:"token"`
	Createtime int64  `json:"createtime"`
}

type HandlerFunc func(token *Token) error

func WithExpired(expired int64) func(token *Token) error {
	return func(token *Token) error {
		if token.Createtime+expired < time.Now().Unix() {
			return fmt.Errorf("token failure")
		}
		return nil
	}
}
