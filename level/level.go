package level

type Level interface {
	Type() string
}

type CustomLevel struct {
	levelType string
}

func NewCustomLevel(_type string) *CustomLevel {
	return &CustomLevel{levelType: _type}
}

func (c *CustomLevel) Type() string {
	return c.levelType
}
