package cli

type Flag interface {
	GetName() string
	GetUsage() string
	Setup(c *Context) FlagGet
}

type FlagGet func() interface{}
