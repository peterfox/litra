package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		// lists out the commands
		fmt.Println("Usage: litra <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  on")
		fmt.Println("  off")
		fmt.Println()
	}

	if err := root(os.Args[1:]); err != nil {
		flag.Usage()

		fmt.Println(err)
		os.Exit(1)
	}
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func root(args []string) error {

	if len(args) < 1 {
		return errors.New("You must pass a command")
	}

	cmds := []Runner{
		NewOnCommand(),
		NewOffCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown command: %s", subcommand)
}
