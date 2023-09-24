package code

// script
const (
	ErrResourceNotFound = iota + 100000
	ErrResourceMustID
	ErrResourceArgs
)

// script
const (
	ScriptParseFailed = iota + 101000
	ScriptNotFound
	ScriptRuntimeNotFound
	ScriptAlreadyEnable
	ScriptAlreadyDisable
	ScriptStateError
	ScriptRunStateError

	StorageNameNotFound
)
