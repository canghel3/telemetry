package main

import (
	"fmt"
	"telemetry/log"
)

type CustomDriver struct {
	msg string
}

func (cd *CustomDriver) Write(p []byte) error {
	cd.msg = string(p)
	return nil
}

func main() {
	log.Stdout().Error().Write([]byte("HELLO"))

	stdout := log.Stdout()

	stdout.Error().Write([]byte("HELLO"))
	stdout.Info().Write([]byte("WORLD"))

	toFile := log.File("./telemetry.log")
	tx := toFile.BeginTx()
	tx.Append(toFile.Error().Append([]byte("hallelujah")))
	tx.Append(toFile.Info().Append([]byte("marcele, la covrigarie!")))
	tx.Commit()

	var customDriver CustomDriver
	customDriver.Write([]byte("ยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩"))
	fmt.Println(customDriver.msg)
}
