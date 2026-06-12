package main

import (
	"errors"
	"fmt"
	"os"

	"litra/driver"
)

func usage() {
	fmt.Println("Usage: litra <command>")
	fmt.Println()
	fmt.Println("Commands:")
	for _, cmd := range commands {
		fmt.Printf("  %-25s %s\n", cmd.usage, cmd.summary)
	}
	fmt.Println()
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		usage()
		return errors.New("you must pass a command")
	}

	for _, cmd := range commands {
		if cmd.name != args[0] {
			continue
		}

		ld, err := driver.New()
		if err != nil {
			return fmt.Errorf("device not found: %w", err)
		}
		defer ld.Close()

		return cmd.run(ld, args[1:])
	}

	usage()
	return fmt.Errorf("unknown command: %s", args[0])
}
