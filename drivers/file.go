package drivers

import (
	"errors"
	"fmt"
	"os"
)

type FileDriver struct {
	name string
	file *os.File
}

func NewFileDriver(name string) *FileDriver {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %s: %s", name, err.Error())
	}

	return &FileDriver{
		name: name,
		file: file,
	}
}

// Log always appends for files.
func (f *FileDriver) Write(p []byte) (int, error) {
	if f.file == nil {
		return 0, errors.New("file not open")
	}

	return f.file.Write(p)
}
