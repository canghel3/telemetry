package drivers

// OutputDriver is identical to io.Writer
type OutputDriver interface {
	Write([]byte) (int, error)
}
