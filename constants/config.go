package constants

import "runtime"

// DefaultConfigContent is the default config config.
const DefaultConfigContent = `concurrency: 0
log_file: ~/.qscamel/qscamel.log
log_level: info
pid_file: ~/.qscamel/qscamel.pid
database_file: ~/.qscamel/db
`

var (
	// DefaultConcurrency is default num of objects being migrated concurrently.
	DefaultConcurrency = runtime.NumCPU() * 10
	// GoldenRatio is the most beautiful number on earth.
	// ref: https://en.wikipedia.org/wiki/Golden_ratio
	GoldenRatio = 0.618
	// MB represent 1024*1024 Byte.
	MB int64 = 1024 * 1024
)

// Path store all path related constants.
const (
	Path         = "~/.qscamel"
	ConfigPath   = Path + "/qscamel.yaml"
	DatabasePath = Path + "/db"
	LogPath      = Path + "/qscamel.log"
	PIDPath      = Path + "/qscamel.pid"
)
