package api

import (
	"context"

	"github.com/codfrm/cago/server/mux"
	"github.com/scriptscat/cloudcat/internal/controller/scripts_ctr"
	"github.com/scriptscat/cloudcat/internal/repository/cookie_repo"
	"github.com/scriptscat/cloudcat/internal/repository/script_repo"
	"github.com/scriptscat/cloudcat/internal/repository/value_repo"
	"github.com/scriptscat/cloudcat/internal/service/scripts_svc"
)

// Router 路由表
// @title    云猫 API 文档
// @version  1.0.0
// @BasePath /api/v1
func Router(ctx context.Context, root *mux.Router) error {
	r := root.Group("/api/v1")

	script_repo.RegisterScript(script_repo.NewScript())
	value_repo.RegisterValue(value_repo.NewValue())
	cookie_repo.RegisterCookie(cookie_repo.NewCookie())

	_, err := scripts_svc.NewScript(ctx)
	if err != nil {
		return err
	}
	{
		script := scripts_ctr.NewScripts()
		r.Bind(
			script.List,
			script.Install,
			script.Update,
			script.Get,
			script.Delete,
		)
	}
	{
		value := scripts_ctr.NewValue()
		r.Bind(
			value.ValueList,
		)
	}
	{
		cookie := scripts_ctr.NewCookie()
		r.Bind(
			cookie.CookieList,
		)
	}

	return nil
}
