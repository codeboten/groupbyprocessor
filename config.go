package groupbyprocessor

type Config struct {
	// Deduplicate iterates thorugh logs and keeps a map of
	// log bodies to remove duplicates from the processed payload
	Deduplicate bool `mapstructure:"deduplicate"`
	// Flatten configures the processor to combine similar logs but
	// combines their attributes
	Flatten bool `mapstructure:"flatten"`
	// MaxLogsBuffered restricts the size of cache to store
	// logs in. Defaults to 64
	MaxLogsBuffered int64 `mapstructure:"max_logs_buffered"`
}
