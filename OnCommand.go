package main

import (
	"flag"
	"fmt"
	"litra/driver"
)

func NewOnCommand() *OnCommand {
	gc := &OnCommand{
		fs: flag.NewFlagSet("on", flag.ContinueOnError),
	}

	return gc
}

type OnCommand struct {
	fs *flag.FlagSet

	name string
}

func (g *OnCommand) Name() string {
	return g.fs.Name()
}

func (g *OnCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *OnCommand) Run() error {
	ld, err := driver.New()

	if err != nil {
		fmt.Println("Device not found")
		return nil
	}

	ld.TurnOn()

	ld.Close()
	return nil
}
