package env

var (
	DbPath              string
	BackupDbPath        string
	BackupCronSchedule  string
	LogLevel            int
	Port                string
	SessionTokenLength  int
	SessionTokenCharset string
	CSRFTokenLength     int
	CSRFTokenCharset    string
	PageSizeMax         int64
	PageSizeDefault     int64
)

func SetConstants() {
	loadOptionalEnvFile()

	DbPath = MustGetString("DB_PATH", "data/live/db/live.db")
	BackupDbPath = MustGetString("BACKUP_DB_PATH", "data/backup/db/backup.db")
	BackupCronSchedule = MustGetString("BACKUP_CRON_SCHEDULE", "0 0 * * *")
	LogLevel = MustGetInt("LOG_LEVEL", 0)
	Port = MustGetString("PORT", "8080")
	SessionTokenLength = MustGetInt("SESSION_TOKEN_LENGTH", 32)
	SessionTokenCharset = MustGetString("SESSION_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	CSRFTokenLength = MustGetInt("CSRF_TOKEN_LENGTH", 32)
	CSRFTokenCharset = MustGetString("CSRF_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	PageSizeMax = MustGetInt64("PAGE_SIZE_MAX", 100)
	PageSizeDefault = MustGetInt64("PAGE_SIZE_DEFAULT", 10)
}
