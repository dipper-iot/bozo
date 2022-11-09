package cli

import (
	"context"
)

type App struct {
	Name        string
	Description string
	Flags       []Flag
	Commands    []*Command
	Before      ActionBeforeFunc
	Action      ActionFunc
	After       ActionFunc
	Context     context.Context
}

type ActionFunc func(c *Context) error
type ActionBeforeFunc func(c *ContextBefore) error

func (a *App) AddCommand(command *Command) {
	a.Commands = append(a.Commands, command)
}

func (a *App) Run(args []string) error {
	appName, argsResult, flagResult, err := Parse(a, args)
	if err != nil {
		return err
	}

	c := &Context{
		Context:     a.Context,
		args:        argsResult,
		cliApp:      a,
		dataVal:     map[string]interface{}{},
		Run:         true,
		dataFlagVal: flagResult,
	}

	flagCall := make(map[string]FlagGet)
	for _, flag := range a.Flags {
		if flag == nil {
			continue
		}
		flagCall[flag.GetName()] = flag.Setup(c)
	}
	c.flagVal = flagCall

	if a.Before != nil {
		err = a.Before(&ContextBefore{
			Context: c,
		})
		if err != nil {
			return err
		}
	}

	commandName, okCommand := c.GetArg(0)

	if c.IsHelp() && !okCommand {
		showHelp(a, appName)
		return nil
	}

	if c.IsHelp() && okCommand {
		for _, command := range a.Commands {
			if command.Name == commandName {
				showHelpCommand(a, command)
				return nil
			}
		}
	}

	if okCommand {
		for _, command := range a.Commands {
			if command.Name == commandName {
				cc := c.Clone()
				cc.args = c.args[1:]
				err = a.RunCommand(command, cc)
				if err != nil {
					return err
				}
			}
		}
	}

	if c.IsHelp() {
		return nil
	}

	err = a.Action(c)
	if err != nil {
		return err
	}

	err = a.After(c)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) RunCommand(command *Command, c *Context) error {

	commandName, okCommand := c.GetArg(0)

	if c.IsHelp() && !okCommand {
		showHelpCommand(a, command)
		return nil
	}

	if okCommand {
		for _, sub := range command.SubCommands {
			if sub.Name == commandName {
				cc := c.Clone()
				cc.args = c.args[1:]
				err := a.RunCommand(sub, cc)
				if err != nil {
					return err
				}
			}
		}
	}
	if c.IsHelp() {
		return nil
	}

	err := command.Action(c)
	if err != nil {
		return err
	}

	return nil
}
