package drivers

import "os"

type Stdout struct {
}

func ToStdout() *Stdout {
	return &Stdout{}
}

func (s *Stdout) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}
