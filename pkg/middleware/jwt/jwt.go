package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const Userinfo = "userinfo"
const JwtToken = "jwt_token"

func Jwt(jwtToken []byte, enforce bool, handlers ...HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth, _ := ctx.Cookie("auth")
		if auth == "" {
			auth = ctx.GetHeader("Authorization")
			if auth == "" {
				auth = ctx.PostForm("auth")
			} else {
				auths := strings.Split(auth, " ")
				if len(auths) != 2 {
					ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code": 1000, "msg": "auth string is empty",
					})
					return
				} else {
					auth = auths[1]
				}
			}
		}
		if auth == "" {
			if enforce {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code": 1001, "msg": "auth string is empty",
				})
			}
			return
		}

		token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			for _, v := range handlers {
				if err := v(token); err != nil {
					return nil, err
				}
			}
			return jwtToken, nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 1002, "msg": err.Error(),
			})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set(Userinfo, claims)
			ctx.Set(JwtToken, token)
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 1003, "msg": "token is wrong",
			})
		}

	}
}

func GenJwt(jwtToken []byte, data jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	token.Header["time"] = time.Now().Unix()
	tokenString, err := token.SignedString(jwtToken)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
