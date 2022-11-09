package cli

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {
	app := &App{
		Name: "app-test",
		Flags: []Flag{
			StringFlag{
				Name:         "test",
				Aliases:      []string{"t"},
				Env:          []string{"TEST_ENV"},
				DefaultValue: "bar",
				Usage:        "Flag Test",
			},
		},
		Commands: []*Command{
			&Command{
				Name: "command1",
				Before: func(c *Context) error {
					fmt.Println("AddCommand Before")
					return nil
				},
				Description: "Command 1",
				Action: func(c *Context) error {
					fmt.Println("AddCommand Action")
					return nil
				},
				SubCommands: []*Command{
					{
						Name:        "sub1",
						Description: "sub command 1",
					},
				},
			},
		},
	}

	err := app.Run(strings.Split("/app -t command1 -h", " "))
	if err != nil {
		log.Fatalln(err)
	}

}

func TestApp2(t *testing.T) {
	app := &App{
		Name: "app-test-2",
		Flags: []Flag{
			StringFlag{
				Name:         "test",
				Aliases:      []string{"t"},
				Env:          []string{"TEST_ENV"},
				DefaultValue: "bar",
				Usage:        "Flag Test",
			},
		},
		Before: func(c *ContextBefore) error {

			c.AddCommand(&Command{
				Name: "command1",
				Before: func(c *Context) error {
					fmt.Println("AddCommand Before")
					return nil
				},
				Description: "Command 1",
				Action: func(c *Context) error {
					fmt.Println("AddCommand Action")
					return nil
				},
				SubCommands: []*Command{
					{
						Name:        "sub1",
						Description: "sub command 1",
					},
				},
			})

			return nil
		},
	}

	err := app.Run(strings.Split("/app command1 -h", " "))
	if err != nil {
		log.Fatalln(err)
	}

}
