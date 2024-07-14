package level

type CustomLevel struct {
	levelType string
}

func Custom(_type string) *CustomLevel {
	return &CustomLevel{levelType: _type}
}

func (c *CustomLevel) Type() string {
	return c.levelType
}
