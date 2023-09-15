package main

import (
	"flag"

	"github.com/derickr/go-litra-driver"
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
	ld, err := litra.New()

	if err != nil {
		panic(err)
	}

	ld.TurnOn()

	ld.Close()
	return nil
}
