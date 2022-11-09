package cli

import (
	"flag"
	"os"
	"strings"
)

type BoolFlag struct {
	Name         string
	Aliases      []string
	Usage        string
	Env          []string
	DefaultValue bool
}

func (s BoolFlag) Setup(c *Context) FlagGet {
	var val bool
	var valAliases = make(map[string]*bool)
	flag.BoolVar(&val, s.Name, s.DefaultValue, s.Usage)
	for _, alias := range s.Aliases {
		flag.BoolVar(valAliases[alias], s.Name, s.DefaultValue, s.Usage)
	}

	return func() interface{} {
		for _, env := range s.Env {
			val, ok := os.LookupEnv(env)
			if ok {
				return covertBool(val)
			}
		}

		val, ok := c.GetFlag(s.Name)
		if ok {
			return covertBool(val)
		}

		for _, name := range s.Aliases {
			val, ok := c.GetFlag(name)
			if ok {
				return covertBool(val)
			}
		}

		return s.DefaultValue
	}
}

func covertBool(val string) bool {
	val = strings.ToLower(val)
	if val == "" || val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	return false
}

func (s BoolFlag) GetName() string {
	return s.Name
}

func (s BoolFlag) GetUsage() string {
	return s.Usage
}
