Cli
=========

[![Build Status](https://travis-ci.org/ionutmilica/cli.svg)](https://travis-ci.org/ionutmilica/cli)
[![Coverage Status](https://coveralls.io/repos/ionutmilica/cli/badge.svg?branch=master&service=github)](https://coveralls.io/github/ionutmilica/cli?branch=master)

This library is made for those who want simplicity when dealing with cli applications.
I'm making this because I have a need for it on my GO web framework, VUA but also for learning purposes.


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

func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:      "Build",
		Signature: "build {what}",
		Action: func() {
			println("Build command!")
		},
	}
}

func ClearCommand() *cli.Command {
	return &cli.Command{
		Name:      "Clear",
		Signature: "clear",
		Action: func() {
			println("Clear command!")
		},
	}
}

````

This project is under development so it's not production ready.


License
----

MIT