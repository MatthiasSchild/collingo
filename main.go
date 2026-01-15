package main

import (
	"collingo/commands"
	"collingo/console"
	"os"
)

func main() {
	err := commands.RootCmd.Execute()
	if err != nil {
		console.Error(err)
		os.Exit(1)
	}
}
