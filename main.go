package main

import (
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/log"
	"github.com/spf13/viper"
	"os"
)

type CustomDriver struct {
	msg string
}

func (cd *CustomDriver) Write(p []byte) error {
	cd.msg = string(p)
	return nil
}

func init() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("TELEMETRY")
	v.AddConfigPath("./")
	v.SetConfigName("telemetry")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		log.Stdout().Error().Logf("failed to read config: %s", err.Error())
	}

	err = v.Unmarshal(&config.PkgConfiguration)
	if err != nil {
		log.Stdout().Error().Logf("failed to unmarshal config: %s", err.Error())
	}
}

// formatting is very basic, but not vital for proper functionality
// solve TODOs
func main() {
	const logfile = "./telemetry.log"
	const configFile = "./config.json"
	log.Stdout().Error().Log("HELLO")

	stdout := log.Stdout()

	stdout.Info().Logf("a formatted log %s", "OI BILLY")
	stdout.Info().Metadata(map[any]any{"something": "clean"}).Log("salutare")

	stdout.Error().Log("HELLO")
	stdout.Info().Log("WORLD")

	//toFile := log.File(logfile)
	//tx := log.BeginTx()
	//tx.Append(toFile.Error().Msg([]byte("something is going on")))
	//tx.Append(toFile.Info().Msg([]byte("marcele, la covrigarie!")))
	//tx.Append(log.Stdout().Msg([]byte("TO STDOUT!")))
	//tx.Log()
	//
	//stdoutWithSettings := log.Stdout()
	//stdoutWithSettings.Settings("./overwriter.json")
	//stdoutWithSettings.Log([]byte("inghetata de fistic"))

	os.WriteFile(logfile, nil, 0644)
}
