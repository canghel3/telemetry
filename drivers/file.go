package drivers

import "os"

type File struct {
	name string
}

func ToFileWithName(name string) *File {
	return &File{name}
}

// Log always appends for files.
func (f *File) Log(p []byte) error {
	file, err := os.OpenFile(f.name, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(p)
	return err
}
