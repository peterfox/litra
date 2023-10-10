package main

import (
	"flag"
	"litra/driver"
)

func NewToggleCommand() *ToggleCommand {
	gc := &ToggleCommand{
		fs: flag.NewFlagSet("toggle", flag.ContinueOnError),
	}

	return gc
}

type ToggleCommand struct {
	fs *flag.FlagSet

	name string
}

func (g *ToggleCommand) Name() string {
	return g.fs.Name()
}

func (g *ToggleCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *ToggleCommand) Run() error {
	ld, err := driver.New()

	if err != nil {
		panic(err)
	}

	state, _ := ld.IsOn()

	if state {
		ld.TurnOff()
	} else {
		ld.TurnOn()
	}

	ld.Close()
	return nil
}
