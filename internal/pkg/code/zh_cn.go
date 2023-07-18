package code

import "github.com/codfrm/cago/pkg/i18n"

func init() {
	i18n.Register(i18n.DefaultLang, zhCN)
}

var zhCN = map[int]string{
	ErrResourceNotFound: "资源不存在",
	ErrResourceMustID:   "必须输入资源id",
	ErrResourceArgs:     "参数错误",

	ScriptParseFailed:     "脚本解析失败",
	ScriptNotFound:        "脚本不存在",
	ScriptRuntimeNotFound: "脚本运行时不存在",
	ScriptAlreadyEnable:   "脚本已经启用",
	ScriptAlreadyDisable:  "脚本已经禁用",
	ScriptStateError:      "脚本状态错误",
	ScriptRunStateError:   "脚本运行状态错误",
}
