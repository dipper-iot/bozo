package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type IntFlag struct {
	Name         string
	Aliases      []string
	Usage        string
	Env          []string
	DefaultValue int
}

func (s IntFlag) Setup(c *Context) FlagGet {
	return func() interface{} {
		for _, env := range s.Env {
			val, ok := os.LookupEnv(env)
			if ok {
				return covertInt(val)
			}
		}

		val, ok := c.GetFlag(s.Name)
		if ok {
			return covertInt(val)
		}

		for _, name := range s.Aliases {
			val, ok := c.GetFlag(name)
			if ok {
				return covertInt(val)
			}
		}
		return s.DefaultValue
	}
}

func covertInt(val string) int {
	valN, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int(valN)
}

func (s IntFlag) GetName() string {
	return s.Name
}

func (s IntFlag) GetUsage() string {
	envStr := strings.Join(s.Env, " ")
	if len(envStr) > 0 {
		envStr = fmt.Sprintf(" [%s]", envStr)
	}
	aliases := ""
	for _, alias := range s.Aliases {
		aliases = fmt.Sprintf("%s, --%s", aliases, alias)
	}
	defaultVal := ""
	if s.DefaultValue > 0 {
		defaultVal = fmt.Sprintf(" (Default: %d)", s.DefaultValue)
	}
	return fmt.Sprintf("-%s%s - %s%s%s", s.Name, aliases, s.Usage, envStr, defaultVal)
}
