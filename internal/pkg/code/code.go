package code

// auth
const (
	TokenIsEmpty = iota + 100000
	TokenIsInvalid
	TokenIsExpired
	TokenNotFound
)

// script
const (
	ErrResourceNotFound = iota + 101000
	ErrResourceMustID
	ErrResourceArgs
)

// script
const (
	ScriptParseFailed = iota + 102000
	ScriptNotFound
	ScriptRuntimeNotFound
	ScriptAlreadyEnable
	ScriptAlreadyDisable
	ScriptStateError
	ScriptRunStateError

	StorageNameNotFound
)
