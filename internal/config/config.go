package config

type RunConfig struct {
	Address               string `env:"ADDRESS"`
	DefaultPollInterval   int    `env:"REPORT_INTERVAL"`
	DefaultReportInterval int    `env:"POLL_INTERVAL"`
}
