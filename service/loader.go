package service

import (
	"github.com/dipper-iot/bozo/cli"
	"sort"
)

type ILoader interface {
	Name() string
	Flags() []cli.Flag
	Priority() int
	//Reference() []string
	Start(o *Options, c *cli.Context) error
	//ReferenceWith(name string, object interface{})
	Stop() error
}

type OptionsLoader struct {
	Name string
}
type OptionLoader func(o *OptionsLoader)

func LoaderName(name string) OptionLoader {
	return func(o *OptionsLoader) {
		o.Name = name
	}
}

func sortLoaders(loaders []ILoader) []ILoader {
	sort.Slice(loaders, func(i, j int) bool {
		return loaders[i].Priority() > loaders[j].Priority()
	})

	return loaders
}

func runLoader(loaders []ILoader, o *Options, c *cli.Context, start bool) error {
	for _, loader := range loaders {
		if start {
			err := loader.Start(o, c)
			if err != nil {
				return err
			}
		} else {
			err := loader.Stop()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runLoaderFlag(loaders []ILoader, o *Options) {
	for _, loader := range loaders {
		o.Flags(loader.Flags())
	}
}
