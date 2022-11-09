package config

import (
	"github.com/dipper-iot/bozo/cli"
	"github.com/dipper-iot/bozo/config/source"
	"github.com/dipper-iot/bozo/config/source/consul"
	"github.com/dipper-iot/bozo/config/source/file"
	"github.com/dipper-iot/bozo/logger"
	"github.com/dipper-iot/bozo/service"
)

type ConfigLoader struct {
	defaultName string
	opts        []source.Option
}

func NewConfigLoader(defaultName string, opts ...source.Option) *ConfigLoader {
	return &ConfigLoader{
		defaultName: defaultName,
		opts:        opts,
	}
}

func (c ConfigLoader) Name() string {
	return "config"
}

func (c ConfigLoader) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:         "config",
			Usage:        "provider config: (file|consul)",
			DefaultValue: c.defaultName,
			Env:          []string{"CONF_TYPE"},
			Aliases:      []string{"c"},
		},
	}
}

func (c ConfigLoader) Priority() int {
	return 0
}

func (c ConfigLoader) Start(o *service.Options, ci *cli.Context) error {
	configType := ci.String("config")
	var sourceData source.Source

	switch configType {
	case "consul":
		sourceData = consul.NewSource(c.opts...)
		break
	default:
		sourceData = file.NewSource(c.opts...)
		break
	}

	configData, err := NewConfig(
		WithSource(sourceData),
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	DefaultConfig = configData
	return nil
}

func (c ConfigLoader) Stop() error {

	return nil
}
