package drivers

type Generic interface {
	Write([]byte) error
}
