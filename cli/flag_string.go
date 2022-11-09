package cli

import (
	"fmt"
	"os"
	"strings"
)

type StringFlag struct {
	Name         string
	Aliases      []string
	Usage        string
	Env          []string
	DefaultValue string
}

func (s StringFlag) Setup(c *Context) FlagGet {

	return func() interface{} {
		for _, env := range s.Env {
			val, ok := os.LookupEnv(env)
			if ok {
				return val
			}
		}

		val, ok := c.GetFlag(s.Name)
		if ok {
			return val
		}

		for _, name := range s.Aliases {
			val, ok := c.GetFlag(name)
			if ok {
				return val
			}
		}

		return s.DefaultValue
	}
}

func (s StringFlag) GetName() string {
	return s.Name
}

func (s StringFlag) GetUsage() string {
	envStr := strings.Join(s.Env, " ")
	if len(envStr) > 0 {
		envStr = fmt.Sprintf(" [%s]", envStr)
	}
	aliases := ""
	for _, alias := range s.Aliases {
		aliases = fmt.Sprintf("%s, --%s", aliases, alias)
	}
	defaultVal := ""
	if len(s.DefaultValue) > 0 {
		defaultVal = fmt.Sprintf(" (Default: %s)", s.DefaultValue)
	}
	return fmt.Sprintf("-%s%s - %s%s%s", s.Name, aliases, s.Usage, envStr, defaultVal)
}
