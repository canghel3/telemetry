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

	var customDriver CustomDriver
	customDriver.Write([]byte("ยัญชนะ(⟨б⟩, ⟨в⟩, ⟨г⟩"))
	fmt.Println(customDriver.msg)
}
