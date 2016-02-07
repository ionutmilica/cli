Cli
=========

[![GoDoc](https://godoc.org/github.com/codegangsta/cli?status.svg)](https://godoc.org/github.com/ionutmilica/cli)
[![Build Status](https://travis-ci.org/ionutmilica/cli.svg)](https://travis-ci.org/ionutmilica/cli)
[![Coverage Status](https://coveralls.io/repos/ionutmilica/cli/badge.svg?branch=master&service=github)](https://coveralls.io/github/ionutmilica/cli?branch=master)

I wanted to make a simple cli app for Vua framework, but I've wanted something simple like the Console from
Laravel so I've made this library.
It provides very simple, but powerful syntax for creating commands for your CLI application.

Example:
```go
package main

import (
	"github.com/ionutmilica/cli"
	"os"
)

func main() {
	app := cli.New()
	app.AddCommand(BuildCommand)
	app.AddCommand(ClearCommand)
	app.Run(os.Args)
}

func BuildCommand(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:        "Build",
		Signature:   "{input} {output=out.bin} {--verbose}",
		Description: "Build this project",
		Action: func(ctx *cli.Context) {
			println("Build command!")
		},
	}
}

func ClearCommand(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:        "Clear",
		Signature:   "",
		Description: "Clears something from the project",
		Action: func(ctx *cli.Context) {
			println("Clear command!")
		},
	}
}
````

This project is under development so it's not production ready.

Todo List
----
- [x] Required argument/ value for long option, i.e {file}
- [x] Optional argument / value for long option, i.e {file?}
- [x] Array argument, i.e {files=*}
- [x] Description for argument, i.e {file : This argument accept a string}
- [ ] Options, i.e {-q}
- [ ] Array value for option
- [ ] Argument default value , i.e {user=johnny}
- [ ] Long Option default value, i.e {--queue=redis}
- [ ] Option alias, i.e {-q|queue}
- [ ] Subcommands
- [ ] Console helpers: confirm, input, table, secret, ask, text color
- [ ] Autocomplete

License
----

MIT