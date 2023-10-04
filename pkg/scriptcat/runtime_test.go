package scriptcat_test

import (
	"context"
	"testing"

	"github.com/codfrm/cago/pkg/logger"
	scriptcat2 "github.com/scriptscat/cloudcat/pkg/scriptcat"
	"github.com/scriptscat/cloudcat/pkg/scriptcat/plugin/window"
	_ "github.com/scriptscat/cloudcat/pkg/scriptcat/plugin/window"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestScriptCat(t *testing.T) {
	ctx := context.Background()
	logger.SetLogger(zap.L())
	r := scriptcat2.NewRuntime(logger.NewCtxLogger(logger.Default()), []scriptcat2.Plugin{
		window.NewBrowserPlugin(),
	})
	script := &scriptcat2.Script{
		ID: "1",
		Code: "return new Promise(resolve=>{" +
			"resolve('ok')" +
			"})",
		Metadata: nil,
	}
	_, err := r.Run(ctx, script)
	assert.Nil(t, err)
}
