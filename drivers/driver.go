package drivers

type OutputDriver interface {
	Write([]byte) error
}
