# Litra

Just a simple utility to work with the Logitech Litra Glow light.

## Install

### Brew

You can install via brew.

```sh
brew tap peterfox/litra
brew install litra
```

## Usage

### Basic controls

```shell
$ litra on
$ litra off
$ litra toggle
$ litra brightness 50      # set brightness 0-100 (percent)
$ litra temperature 4000   # set colour temperature 2700-6500 (kelvin)
```

### Relative adjustments

Increase or decrease brightness and colour temperature relative to the current device state. An optional step value can be provided; if omitted a sensible default is used.

```shell
$ litra brightness-up         # increase brightness by 10%
$ litra brightness-up 20      # increase brightness by 20%
$ litra brightness-down       # decrease brightness by 10%
$ litra brightness-down 5     # decrease brightness by 5%

$ litra temperature-up        # increase colour temperature by 200K
$ litra temperature-up 500    # increase colour temperature by 500K
$ litra temperature-down      # decrease colour temperature by 200K
$ litra temperature-down 100  # decrease colour temperature by 100K
```

Values are clamped to the valid range (brightness: 0–100, temperature: 2700–6500) so there is no need to track where you are.

### Profiles

Profiles let you save the current brightness and temperature under a name and restore them later. They are stored in `~/.litra`.

```shell
# Save the current light state as a profile
$ litra profile save work

# Load a saved profile
$ litra profile load work

# List all saved profiles
$ litra profile list
```

Example `~/.litra` file:

```json
{
  "profiles": {
    "work":    { "brightness": 70, "temperature": 4500 },
    "evening": { "brightness": 30, "temperature": 2800 },
    "video":   { "brightness": 85, "temperature": 5500 }
  }
}
```

The file can also be edited by hand — any valid JSON matching the structure above will be picked up by the CLI.

## Build

```shell
go build -v .
```

## License

MIT
