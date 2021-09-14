package v1

import (
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt"
	"github.com/scriptscat/cloudcat/pkg/middleware/jwt"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

func userId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(jwt.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(jwt2.MapClaims)["uid"].(string)), true
}

func isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(jwt.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(jwt2.MapClaims)["uid"].(string)), false
}

func jwttoken(ctx *gin.Context) (jwt2.MapClaims, *jwt2.Token, bool) {
	u, ok := ctx.Get(jwt.Userinfo)
	if !ok {
		return nil, nil, false
	}
	t, ok := ctx.Get(jwt.JwtToken)
	if !ok {
		return nil, nil, false
	}
	return u.(jwt2.MapClaims), t.(*jwt2.Token), true
}
