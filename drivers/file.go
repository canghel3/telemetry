package drivers

import "os"

type File struct {
	name string
}

func ToFileWithName(name string) *File {
	return &File{name}
}

// Write always appends for files.
func (f *File) Write(p []byte) error {
	file, err := os.OpenFile(f.name, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(p)
	return err
}
