package env

var (
	DbPath              = MustGetString("DB_PATH", "data/live/db/data.db")
	LogLevel            = MustGetInt("LOG_LEVEL", 0)
	Port                = MustGetString("PORT", "8080")
	SessionTokenLength  = MustGetInt("SESSION_TOKEN_LENGTH", 16)
	SessionTokenCharset = MustGetString("SESSION_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	CSRFTokenLength     = MustGetInt("CSRF_TOKEN_LENGTH", 16)
	CSRFTokenCharset    = MustGetString("CSRF_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)
