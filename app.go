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

// Register a new command into the system
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

	if len(osArgs) > 1 {
		cmd = osArgs[1]
		args = osArgs[2:]
	}

	if cmd, ok := app.Commands[cmd]; ok {
		cmd.parse()

		matcher := newMatcher(args, cmd.Flags)

		if err := matcher.match(); err != nil {
			fmt.Println(err.Error())
			return
		}
		cmd.Action(matcher.ctx)
		return
	}

	fmt.Printf("Command `%s` was not found!", cmd)
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
