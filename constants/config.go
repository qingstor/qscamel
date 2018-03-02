package constants

// DefaultConfigContent is the default config config.
const DefaultConfigContent = `thread_num: 100
log_file: ~/.qscamel/qscamel.log
log_level: error
database_file: ~/.qscamel/qscamel.db
`

const (
	// DefaultThreadNum is default num of objects being migrated concurrently.
	DefaultThreadNum = 100
)

// Path store all path related constants.
const (
	Path         = "~/.qscamel"
	ConfigPath   = Path + "/qscamel.yaml"
	DatabasePath = Path + "/qscamel.db"
	LogPath      = Path + "/qscamel.log"
)
