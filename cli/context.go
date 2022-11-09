package cli

import "context"

type Context struct {
	cliApp      *App
	Context     context.Context
	dataVal     map[string]interface{}
	dataFlagVal map[string]string
	flagVal     map[string]FlagGet
	Run         bool
	Commands    []string
	args        []string
	App         string
}

func (c Context) Clone() *Context {
	return &Context{
		cliApp:      c.cliApp,
		Context:     c.Context,
		dataVal:     c.dataVal,
		dataFlagVal: c.dataFlagVal,
		flagVal:     c.flagVal,
		Run:         c.Run,
		Commands:    c.Commands,
		args:        c.args,
		App:         c.App,
	}
}

type ContextBefore struct {
	*Context
}

func (c ContextBefore) AddCommand(command *Command) {
	c.cliApp.AddCommand(command)
}

func (c *Context) GetFlag(name string) (string, bool) {
	val, ok := c.dataFlagVal[name]
	return val, ok
}

func (c *Context) Get(name string) (interface{}, bool) {
	val, ok := c.dataVal[name]
	return val, ok
}

func (c *Context) Set(name string, val interface{}) {
	c.dataVal[name] = val
}

func (c *Context) Args() []string {
	return c.args
}

func (c *Context) GetArg(i int) (string, bool) {
	if len(c.args) <= i {
		return "", false
	}
	return c.args[i], true
}

func (c Context) String(name string) string {
	call, ok := c.flagVal[name]
	if !ok {
		return ""
	}
	val := call()
	return val.(string)
}

func (c Context) Int(name string) int {
	call, ok := c.flagVal[name]
	if !ok {
		return 0
	}
	val := call()
	return val.(int)
}

func (c Context) Bool(name string) bool {
	call, ok := c.flagVal[name]
	if !ok {
		return false
	}
	val := call()
	return val.(bool)
}

func (c Context) Has(name string) bool {
	_, ok := c.flagVal[name]
	return ok
}

func (c Context) IsHelp() bool {
	_, ok := c.dataFlagVal["h"]
	if ok {
		return true
	}
	_, ok = c.dataFlagVal["help"]
	return ok
}
