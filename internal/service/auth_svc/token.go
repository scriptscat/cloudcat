package auth_svc

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/codfrm/cago/pkg/i18n"
	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	api "github.com/scriptscat/cloudcat/internal/api/auth"
	"github.com/scriptscat/cloudcat/internal/model/entity/token_entity"
	"github.com/scriptscat/cloudcat/internal/pkg/code"
	"github.com/scriptscat/cloudcat/internal/repository/token_repo"
	utils2 "github.com/scriptscat/cloudcat/pkg/utils"
)

type TokenSvc interface {
	// TokenList 获取token列表
	TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error)
	// TokenCreate 创建token
	TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error)
	// TokenDelete 删除token
	TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error)
	// Middleware 验证token中间件
	Middleware() gin.HandlerFunc
}

type tokenSvc struct {
}

var defaultToken = &tokenSvc{}

func Token() TokenSvc {
	return defaultToken
}

// TokenList 获取token列表
func (s *tokenSvc) TokenList(ctx context.Context, req *api.TokenListRequest) (*api.TokenListResponse, error) {
	list, err := token_repo.Token().FindPage(ctx)
	if err != nil {
		return nil, err
	}
	resp := &api.TokenListResponse{
		List: make([]*api.Token, 0),
	}
	for _, v := range list {
		if req.TokenID != "" && strings.HasPrefix(v.ID, req.TokenID) {
			continue
		}
		resp.List = append(resp.List, &api.Token{
			ID:                v.ID,
			Token:             v.Token,
			DataEncryptionKey: v.DataEncryptionKey,
			Createtime:        v.Createtime,
			Updatetime:        v.Createtime,
		})
	}
	return resp, nil
}

// TokenCreate 创建token
func (s *tokenSvc) TokenCreate(ctx context.Context, req *api.TokenCreateRequest) (*api.TokenCreateResponse, error) {
	id := req.TokenID
	secret := utils.RandString(32, utils.Letter)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:  "cloudcat",
		Subject: "auth:" + id,
		ID:      id,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	m := &token_entity.Token{
		ID:                id,
		Token:             tokenString,
		Secret:            secret,
		DataEncryptionKey: utils.RandString(32, utils.Letter),
		Createtime:        time.Now().Unix(),
		Updatetime:        time.Now().Unix(),
	}
	if err := token_repo.Token().Create(ctx, m); err != nil {
		return nil, err
	}
	return &api.TokenCreateResponse{
		Token: &api.Token{
			ID:                m.ID,
			Token:             m.Token,
			DataEncryptionKey: m.DataEncryptionKey,
			Createtime:        m.Createtime,
			Updatetime:        m.Createtime,
		}}, nil
}

// TokenDelete 删除token
func (s *tokenSvc) TokenDelete(ctx context.Context, req *api.TokenDeleteRequest) (*api.TokenDeleteResponse, error) {
	err := token_repo.Token().Delete(ctx, req.TokenID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Middleware 验证token中间件
func (s *tokenSvc) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			httputils.HandleResp(c, i18n.NewUnauthorizedError(c, code.TokenIsEmpty))
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		// 解析token
		var m *token_entity.Token
		_, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, i18n.NewUnauthorizedError(c, code.TokenIsInvalid)
			}
			claims := token.Claims.(*jwt.RegisteredClaims)
			var err error
			m, err = token_repo.Token().Find(c, claims.ID)
			if err != nil {
				return nil, err
			}
			if m == nil {
				return nil, i18n.NewUnauthorizedError(c, code.TokenIsInvalid)
			}
			return []byte(m.Secret), nil
		})
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		// 解密body
		oldBody := c.Request.Body
		r, err := utils2.NewAesDecrypt([]byte(m.DataEncryptionKey), c.Request.Body)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		c.Request.Body = utils2.WarpCloser(r, oldBody)
		c.Request.GetBody = func() (io.ReadCloser, error) {
			return c.Request.Body, nil
		}
		// 包装response 加密body
		warpW, err := newWarpWrite(c.Writer, m.DataEncryptionKey)
		if err != nil {
			httputils.HandleResp(c, err)
			return
		}
		c.Writer = warpW
		defer func() {
			_ = warpW.Close()
		}()
		c.Next()
	}
}

type warpWrite struct {
	gin.ResponseWriter
	done chan struct{}
	pr   *io.PipeReader
	pw   *io.PipeWriter
	aes  io.Reader
}

func newWarpWrite(w gin.ResponseWriter, key string) (*warpWrite, error) {
	pr, pw := io.Pipe()
	aes, err := utils2.NewAesEncrypt([]byte(key), pr)
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		defer func() {
			close(done)
			_ = pr.Close()
		}()
		_, _ = io.Copy(w, aes)
	}()
	return &warpWrite{
		ResponseWriter: w,
		done:           done,
		pw:             pw,
		pr:             pr,
		aes:            aes,
	}, nil
}

func (w *warpWrite) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *warpWrite) Write(p []byte) (int, error) {
	return w.pw.Write(p)
}

func (w *warpWrite) Close() error {
	if err := w.pw.Close(); err != nil {
		return err
	}
	<-w.done
	return nil
}
