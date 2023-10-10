package main

import (
	"flag"

	"litra/driver"
)

func NewOffCommand() *OffCommand {
	gc := &OffCommand{
		fs: flag.NewFlagSet("off", flag.ContinueOnError),
	}

	return gc
}

type OffCommand struct {
	fs *flag.FlagSet

	name string
}

func (g *OffCommand) Name() string {
	return g.fs.Name()
}

func (g *OffCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *OffCommand) Run() error {
	ld, err := driver.New()

	if err != nil {
		panic(err)
	}

	ld.TurnOff()

	ld.Close()
	return nil
}
