package main

import (
	"fmt"
	"strconv"

	"litra/driver"
)

func optionalStepArg(args []string, name string, defaultVal int) (int, error) {
	if len(args) == 0 {
		return defaultVal, nil
	}
	if len(args) != 1 {
		return 0, fmt.Errorf("%s accepts an optional positive step value", name)
	}
	step, err := strconv.Atoi(args[0])
	if err != nil || step <= 0 {
		return 0, fmt.Errorf("%s step must be a positive number", name)
	}
	return step, nil
}

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
	{
		name:    "brightness-up",
		usage:   "brightness-up [step]",
		summary: "increase brightness by step percent (default 10)",
		run: func(ld *driver.LitraDevice, args []string) error {
			step, err := optionalStepArg(args, "brightness-up", 10)
			if err != nil {
				return err
			}
			state, err := ld.GetState()
			if err != nil {
				return err
			}
			return ld.SetBrightness(min(state.Brightness+step, 100))
		},
	},
	{
		name:    "brightness-down",
		usage:   "brightness-down [step]",
		summary: "decrease brightness by step percent (default 10)",
		run: func(ld *driver.LitraDevice, args []string) error {
			step, err := optionalStepArg(args, "brightness-down", 10)
			if err != nil {
				return err
			}
			state, err := ld.GetState()
			if err != nil {
				return err
			}
			return ld.SetBrightness(max(state.Brightness-step, 0))
		},
	},
	{
		name:    "temperature-up",
		usage:   "temperature-up [step]",
		summary: "increase colour temperature by step kelvin (default 200)",
		run: func(ld *driver.LitraDevice, args []string) error {
			step, err := optionalStepArg(args, "temperature-up", 200)
			if err != nil {
				return err
			}
			state, err := ld.GetState()
			if err != nil {
				return err
			}
			return ld.SetTemperature(min(state.Temperature+step, 6500))
		},
	},
	{
		name:    "temperature-down",
		usage:   "temperature-down [step]",
		summary: "decrease colour temperature by step kelvin (default 200)",
		run: func(ld *driver.LitraDevice, args []string) error {
			step, err := optionalStepArg(args, "temperature-down", 200)
			if err != nil {
				return err
			}
			state, err := ld.GetState()
			if err != nil {
				return err
			}
			return ld.SetTemperature(max(state.Temperature-step, 2700))
		},
	},
	{
		name:    "profile",
		usage:   "profile <save <name>|load <name>|list>",
		summary: "manage named light profiles in ~/.litra",
		run: func(ld *driver.LitraDevice, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("profile requires a subcommand: save, load, or list")
			}
			switch args[0] {
			case "list":
				return profileList()
			case "save":
				if len(args) < 2 {
					return fmt.Errorf("profile save requires a name")
				}
				return profileSave(ld, args[1])
			case "load":
				if len(args) < 2 {
					return fmt.Errorf("profile load requires a name")
				}
				return profileLoad(ld, args[1])
			default:
				return fmt.Errorf("unknown profile subcommand: %s (use save, load, or list)", args[0])
			}
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
