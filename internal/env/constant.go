package env

var (
	DbPath   = MustGetString("DB_PATH", "data/live/db/data.db")
	LogLevel = MustGetInt("LOG_LEVEL", 0)
)
