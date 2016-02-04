package cli

import (
	"fmt"
	"strings"
)

type App struct {
	commands map[string]*Command
}

func New() *App {
	return &App{
		commands: make(map[string]*Command, 0),
	}
}

func (app *App) AddCommand(cmdFunc func() *Command) {
	c := cmdFunc()

	app.commands[strings.ToLower(c.Name)] = c
}

func (app *App) Run(osArgs []string) {
	if len(osArgs) == 1 {
		// No args
		println("Those are the cmds!")
		return
	}

	cmd := osArgs[1]
	args := osArgs[2:]

	if cmd, ok := app.commands[cmd]; ok {
		cmd.parse()
		cmd.Action(app.createContext(args, cmd.Flags))
		return
	}

	fmt.Printf("Command `%s` was not found!", cmd)
}

func (app *App) createContext(args []string, flags []*Flag) *Context {
	ctx := &Context{
		Arguments: make(map[string][]string, 0),
		Options:   make(map[string][]string),
	}
	ctx.parse(args, newFlagMgr(flags))

	return ctx
}
