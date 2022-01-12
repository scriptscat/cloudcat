package repository

type ScriptCatInfo struct {
	Version string `json:"version"`
	Notice  string `json:"notice"`
}

type System interface {
	GetScriptCatInfo() (*ScriptCatInfo, error)
}
