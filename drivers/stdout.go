package drivers

import "os"

type StdoutDriver struct {
}

func NewStdoutDriver() *StdoutDriver {
	return &StdoutDriver{}
}

func (s *StdoutDriver) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}
