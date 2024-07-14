package drivers

import "os"

type Stdout struct {
}

func ToStdout() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Log(p []byte) error {
	_, err := os.Stdout.Write(p)
	return err
}
