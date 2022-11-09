package cli

import (
	"fmt"
)

func showHelp(a *App, appRun string) {
	fmt.Printf("Application: %s\n", a.Name)
	fmt.Println("-------------------------------")
	fmt.Println(fmt.Sprintf("%s", appRun))
	if len(a.Description) > 0 {
		fmt.Println()
		fmt.Printf("%s\n", a.Description)
		fmt.Println("-------------------------------")
	}
	if len(a.Flags) > 0 {
		fmt.Println("Flag:")
		for _, flag := range a.Flags {
			fmt.Printf("\t%s\n", flag.GetUsage())
		}
		fmt.Println()
	}

	if len(a.Commands) > 0 {
		fmt.Println("Command:")
		for _, command := range a.Commands {
			fmt.Printf("\t%s - %s", command.Name, command.Description)
		}
		fmt.Println()
	}
}

func showHelpCommand(a *App, command *Command) {
	fmt.Printf("Application: %s\n", a.Name)
	fmt.Println("-------------------------------")
	if len(command.Description) > 0 {
		fmt.Println()
		fmt.Printf("%s\n", command.Description)
		fmt.Println("-------------------------------")
	}
	for _, flag := range command.Flags {
		fmt.Println()
		fmt.Println("Flag:")
		fmt.Printf("\t%s\n", flag.GetUsage())
		fmt.Println("-------------------------------")
	}
	if len(command.SubCommands) > 0 {
		fmt.Println("Command:")
		for _, c := range command.SubCommands {
			fmt.Printf("\t%s - %s", c.Name, c.Description)
		}
		fmt.Println()
	}
}
