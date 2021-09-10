package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type HandlerFunc func(token *jwt.Token) error

func WithExpired(expired int64) func(token *jwt.Token) error {
	return func(token *jwt.Token) error {
		if t, ok := token.Header["time"]; ok {
			if i, ok := t.(float64); !(ok && int64(i)+expired > time.Now().Unix()) {
				return fmt.Errorf("token failure")
			}
		} else {
			return fmt.Errorf("token failure")
		}
		return nil
	}
}
