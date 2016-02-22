package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Cli framework main struct
type App struct {
	Commands map[string]*Command
	Writer   io.Writer
	Reader   io.Reader
}

// Creates a new App struct and adds the null command to it
func New() *App {
	app := &App{
		Commands: make(map[string]*Command, 0),
		Writer:   os.Stdout,
		Reader:   os.Stdin,
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
func (app *App) Run(args []string) {
	var cmd string

	args = args[1:]
	cmd, pos := findFirstArgument(args)

	if pos != -1 {
		args = append(args[:pos], args[pos+1:]...)
	}

	if cmd, ok := app.Commands[cmd]; ok {
		cmd.parse()
		ctx := newContext(app.Reader, app.Writer)
		matcher := newMatcher(ctx, args, cmd.Flags)

		if err := matcher.match(); err != nil {
			fmt.Fprintln(app.Writer, err.Error())
			return
		}
		cmd.Action(matcher.ctx)
		return
	}

	fmt.Fprintf(app.Writer, "Command `%s` was not found!", cmd)
}

// Find the first argument from the os args
func findFirstArgument(args []string) (string, int) {
	for i, arg := range args {
		if len(arg) > 0 && arg[0] != '-' {
			return arg, i
		}
	}
	return "", -1
}

// Command for default app usage
func homeCommand(app *App) *Command {
	return &Command{
		Name: "",
		Action: func(ctx *Context) {
			fmt.Fprintln(app.Writer, "Usage:")
			fmt.Fprintln(app.Writer, "\tapp command [arguments]")
			fmt.Fprintln(app.Writer, "The commands are:")
			for _, cmd := range app.Commands {
				if cmd.Name != "" {
					fmt.Fprintf(app.Writer, "\t%s - %s\n", cmd.Name, cmd.Description)
				}
			}
		},
	}
}
