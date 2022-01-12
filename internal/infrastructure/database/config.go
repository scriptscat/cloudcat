package database

type Config struct {
	Type  string
	Mysql struct {
		Dsn    string
		Prefix string
	}
	Sqlite struct {
		File string
		//User   string
		//Passwd string
	}
}
