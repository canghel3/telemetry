package drivers

import (
	"errors"
	"fmt"
	"os"
)

type File struct {
	name string
	file *os.File
}

func ToFileWithName(name string) *File {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %s: %s", name, err.Error())
	}

	return &File{
		name: name,
		file: file,
	}
}

// Log always appends for files.
func (f *File) Write(p []byte) (int, error) {
	if f.file == nil {
		return 0, errors.New("file not open")
	}

	return f.file.Write(p)
}
