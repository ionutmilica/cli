package cli

import (
	"fmt"
	"strings"
)

// Cli framework main struct
type App struct {
	Commands map[string]*Command
}

// Creates a new App struct and adds the null command to it
func New() *App {
	app := &App{
		Commands: make(map[string]*Command, 0),
	}
	app.AddCommand(homeCommand)

	return app
}

func (app *App) AddCommand(cmdFunc func(*App) *Command) *App {
	c := cmdFunc(app)
	app.Commands[strings.ToLower(c.Name)] = c
	return app
}

// Start the cli framework, based on the os arguments. Those arguments should
// follow the pattern: arg1 file, arg2 argument/option and so on
func (app *App) Run(osArgs []string) {
	var cmd string
	var args []string

	if len(args) > 1 {
		cmd = osArgs[1]
		args = osArgs[2:]
	}

	if cmd, ok := app.Commands[cmd]; ok {
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

func homeCommand(app *App) *Command {
	return &Command{
		Name: "",
		Action: func(ctx *Context) {
			fmt.Println("Usage:")
			fmt.Println("\tapp command [arguments]")
			fmt.Println("The commands are:")
			for _, cmd := range app.Commands {
				if cmd.Name != "" {
					fmt.Printf("\t%s - %s\n", cmd.Name, cmd.Description)
				}
			}
		},
	}
}
