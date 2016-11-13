package settings

const (
	HTTP_TIMEOUT = 3
	CACHE_TIME   = 60
	LOCK_TIME    = 10
	LOCK_FILE    = "/tmp/.libstns_lock"
	WORK_DIR     = "/var/lib/libnss_stns"
	MIN_LIMIT_ID = 100
)

const (
	V2_FORMAT_ERROR = "migrate v2 format error. use stns 0.0.6 higer"
)
