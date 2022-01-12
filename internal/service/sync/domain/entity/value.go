package entity

type SyncValue struct {
	UserID      int64
	DeviceID    int64
	ScriptUUID  string
	StorageName string
	Key         string
	Value       string
}
