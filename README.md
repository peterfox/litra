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

```shell
$ litra on
$ litra off
$ litra toggle
$ litra brightness 50      # 0-100 (percent)
$ litra temperature 4000   # 2700-6500 (kelvin)
```

## Build

```shell
go build -v .
```

## License

MIT
