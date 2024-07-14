package drivers

type OutputDriver interface {
	Log([]byte) error
}
