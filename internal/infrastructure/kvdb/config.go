package kvdb

type Config struct {
	Type  string
	Redis struct {
		Addr   string
		Passwd string
		DB     int
	}
	Sqlite struct {
		File string
		//User   string
		//Passwd string
	}
}
