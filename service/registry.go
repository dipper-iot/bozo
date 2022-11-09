package service

import (
	"github.com/dipper-iot/bozo/cli"
	"github.com/dipper-iot/bozo/registry"
	"github.com/dipper-iot/bozo/registry/consul"
)

type RegisterLoader struct {
	register    registry.Registry
	defaultName string
	service     *registry.Service
	opts        []registry.Option
}

func NewRegisterLoader(defaultName string, opts ...registry.Option) *RegisterLoader {
	return &RegisterLoader{defaultName: defaultName, opts: opts}
}

func (r RegisterLoader) Name() string {
	return "registry"
}

func (r RegisterLoader) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:         "registry",
			Aliases:      []string{"r"},
			Env:          []string{"REGISTRY_TYPE"},
			Usage:        "Registry type",
			DefaultValue: r.defaultName,
		},
	}
}

func (r RegisterLoader) Priority() int {
	return 0
}

func (r RegisterLoader) Start(o *Options, c *cli.Context) error {

	switch r.defaultName {
	case "consul":
		r.register = consul.NewRegistry(r.opts...)
		o.Registry = r.register
		registry.DefaultRegistry = r.register
		break
	default:
		r.register = registry.NewDefaultRegistry(r.opts...)
		o.Registry = r.register
		registry.DefaultRegistry = r.register
		break
	}

	return nil
}

func (r RegisterLoader) Stop() error {

	return nil
}
