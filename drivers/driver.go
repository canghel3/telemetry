package drivers

// OutputDriver adheres to io.Writer
type OutputDriver interface {
	Write([]byte) (int, error)
}
