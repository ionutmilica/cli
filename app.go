package cli

import (
	"fmt"
	"strings"
)

func New() *App {
	return &App{
		commands: make(map[string]*Command, 0),
	}
}

type App struct {
	commands map[string]*Command
}

func (app *App) AddCommand(cmdFunc func() *Command) {
	c := cmdFunc()

	app.commands[strings.ToLower(c.Name)] = c
}

func (app *App) Run(args []string) {
	if len(args) == 1 {
		// No args
		println("Those are the cmds!")
		return
	}

	cmd := args[1]
	//_ := args[2:]

	if cmd, ok := app.commands[cmd]; ok {
		cmd.Action()
		return
	}

	fmt.Printf("Command `%s` was not found!", cmd)
}
