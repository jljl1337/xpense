package env

import "net/http"

var (
	Version = "0.3.2"

	DbPath                string
	DbBusyTimeout         string
	BackupDbPath          string
	BackupCronSchedule    string
	LogLevel              int
	LogHealthCheck        bool
	Port                  string
	CORSOrigins           string
	PasswordBcryptCost    int
	SessionCookieName     string
	SessionCookieHttpOnly bool
	SessionCookieSecure   bool
	SessionTokenLength    int
	SessionTokenCharset   string
	SessionLifetimeMin    int
	PreSessionLifetimeMin int
	CSRFTokenLength       int
	CSRFTokenCharset      string
	PageSizeMax           int64
	PageSizeDefault       int64

	SessionCookieSameSiteMode http.SameSite
)

func MustSetConstants() {
	mustLoadOptionalEnvFile()

	DbPath = MustGetString("DB_PATH", "data/live/db/live.db")
	DbBusyTimeout = MustGetString("DB_BUSY_TIMEOUT", "30000")
	BackupDbPath = MustGetString("BACKUP_DB_PATH", "data/backup/db/backup.db")
	BackupCronSchedule = MustGetString("BACKUP_CRON_SCHEDULE", "0 0 * * *")
	LogLevel = MustGetInt("LOG_LEVEL", 0)
	LogHealthCheck = MustGetBool("LOG_HEALTH_CHECK", false)
	Port = MustGetString("PORT", "8080")
	CORSOrigins = MustGetString("CORS_ORIGINS", "*")
	PasswordBcryptCost = MustGetInt("PASSWORD_BCRYPT_COST", 12)
	SessionCookieName = MustGetString("SESSION_COOKIE_NAME", "xpense_session_token")
	SessionCookieHttpOnly = MustGetBool("SESSION_COOKIE_HTTP_ONLY", true)
	SessionCookieSecure = MustGetBool("SESSION_COOKIE_SECURE", false)
	SessionTokenLength = MustGetInt("SESSION_TOKEN_LENGTH", 32)
	SessionTokenCharset = MustGetString("SESSION_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	SessionLifetimeMin = MustGetInt("SESSION_LIFETIME_MIN", 60*24*7)
	PreSessionLifetimeMin = MustGetInt("PRE_SESSION_LIFETIME_MIN", 15)
	CSRFTokenLength = MustGetInt("CSRF_TOKEN_LENGTH", 32)
	CSRFTokenCharset = MustGetString("CSRF_TOKEN_CHARSET", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	PageSizeMax = MustGetInt64("PAGE_SIZE_MAX", 100)
	PageSizeDefault = MustGetInt64("PAGE_SIZE_DEFAULT", 10)

	sessionCookieSameSite := MustGetString("SESSION_COOKIE_SAME_SITE_MODE", "lax")
	switch sessionCookieSameSite {
	case "lax":
		SessionCookieSameSiteMode = http.SameSiteLaxMode
	case "strict":
		SessionCookieSameSiteMode = http.SameSiteStrictMode
	default:
		SessionCookieSameSiteMode = http.SameSiteNoneMode
	}
}
