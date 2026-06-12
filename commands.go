package main

import (
	"fmt"
	"strconv"

	"litra/driver"
)

type command struct {
	name    string
	usage   string
	summary string
	run     func(ld *driver.LitraDevice, args []string) error
}

var commands = []command{
	{
		name:    "on",
		usage:   "on",
		summary: "turn the light on",
		run: func(ld *driver.LitraDevice, _ []string) error {
			return ld.TurnOn()
		},
	},
	{
		name:    "off",
		usage:   "off",
		summary: "turn the light off",
		run: func(ld *driver.LitraDevice, _ []string) error {
			return ld.TurnOff()
		},
	},
	{
		name:    "toggle",
		usage:   "toggle",
		summary: "toggle the light on or off",
		run: func(ld *driver.LitraDevice, _ []string) error {
			on, err := ld.IsOn()
			if err != nil {
				return err
			}
			if on {
				return ld.TurnOff()
			}
			return ld.TurnOn()
		},
	},
	{
		name:    "brightness",
		usage:   "brightness <0-100>",
		summary: "set the brightness as a percentage",
		run: func(ld *driver.LitraDevice, args []string) error {
			level, err := intArg(args, "brightness", 0, 100)
			if err != nil {
				return err
			}
			return ld.SetBrightness(level)
		},
	},
	{
		name:    "temperature",
		usage:   "temperature <2700-6500>",
		summary: "set the colour temperature in kelvin",
		run: func(ld *driver.LitraDevice, args []string) error {
			temp, err := intArg(args, "temperature", 2700, 6500)
			if err != nil {
				return err
			}
			return ld.SetTemperature(temp)
		},
	},
}

func intArg(args []string, name string, lo, hi int) (int, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("%s requires a single value between %d and %d", name, lo, hi)
	}

	value, err := strconv.Atoi(args[0])
	if err != nil || value < lo || value > hi {
		return 0, fmt.Errorf("%s must be a number between %d and %d", name, lo, hi)
	}

	return value, nil
}
