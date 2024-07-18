package config

type PkgConfig struct {
	Formatting FormattingConfig `mapstructure:"formatting"`
}

type FormattingConfig struct {
	LogConfig LogConfig `mapstructure:"log"`
	TxConfig  TxConfig  `mapstructure:"transaction"`
}

type LogConfig struct {
	FormattingDisabled bool           `mapstructure:"disabled"`
	Timestamp          string         `mapstructure:"timestamp"`
	FieldOrder         map[string]int `mapstructure:"field_order"`
}

type TxConfig struct {
	FormattingDisabled bool           `mapstructure:"disabled"`
	Timestamp          string         `mapstructure:"timestamp"`
	FieldOrder         map[string]int `mapstructure:"field_order"`
}

var PkgConfiguration PkgConfig
