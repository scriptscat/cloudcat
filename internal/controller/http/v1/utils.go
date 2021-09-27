package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/scriptscat/cloudcat/pkg/middleware/token"
	"github.com/scriptscat/cloudcat/pkg/utils"
)

func userId(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(token.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), true
}

func isadmin(ctx *gin.Context) (int64, bool) {
	u, ok := ctx.Get(token.Userinfo)
	if !ok {
		return 0, false
	}
	return utils.StringToInt64(u.(gin.H)["uid"].(string)), false
}

func authtoken(ctx *gin.Context) (*token.Token, bool) {
	t, ok := ctx.Get(token.AuthToken)
	if !ok {
		return nil, false
	}
	return t.(*token.Token), true
}
