package config

type Config struct {
	Formatting FormattingConfig `yaml:"formatting"`
}

type FormattingConfig struct {
	LogConfig LogConfig `mapstructure:"log"`
	TxConfig  TxConfig  `mapstructure:"tx"`
}

type LogConfig struct {
	Enabled    bool           `mapstructure:"enabled"`
	Timestamp  string         `mapstructure:"timestamp"`
	FieldOrder map[string]int `mapstructure:"field_order"`
}

type TxConfig struct {
	Enabled    bool           `mapstructure:"enabled"`
	Timestamp  string         `mapstructure:"timestamp"`
	FieldOrder map[string]int `mapstructure:"field_order"`
}
