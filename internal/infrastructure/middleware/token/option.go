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

type Option func(o *options)

type options struct {
	authFailed       []func(token *Token) error
	tokenHandlerFunc []func(token *Token) error
}

type HandlerFunc func(token *Token) error

func WithExpired(expired int64) func(o *options) {
	return func(o *options) {
		o.tokenHandlerFunc = append(o.tokenHandlerFunc, func(token *Token) error {
			if token.Createtime+expired < time.Now().Unix() {
				return fmt.Errorf("token failure")
			}
			return nil
		})
	}
}

func WithDebug(defUser gin.H) func(o *options) {
	return func(o *options) {
		o.authFailed = append(o.authFailed, func(token *Token) error {
			token.Info = defUser
			token.Createtime = time.Now().Unix()
			return nil
		})
	}
}
