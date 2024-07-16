package drivers

import "os"

type File struct {
	name string
}

func ToFileWithName(name string) *File {
	return &File{name}
}

// Log always appends for files.
func (f *File) Write(p []byte) (int, error) {
	file, err := os.OpenFile(f.name, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(p)
}
