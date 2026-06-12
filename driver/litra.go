package driver

import (
	"fmt"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	vendor  = 0x046d
	product = 0xc900

	reportLength = 20

	minBrightness = 0x14
	maxBrightness = 0xfa

	minTemperature = 2700
	maxTemperature = 6500
)

type LitraDevice struct {
	dev *hid.Device
}

func New() (*LitraDevice, error) {
	if err := hid.Init(); err != nil {
		return nil, fmt.Errorf("initialising hid: %w", err)
	}

	dev, err := hid.OpenFirst(vendor, product)
	if err != nil {
		return nil, fmt.Errorf("opening device: %w", err)
	}

	return &LitraDevice{dev: dev}, nil
}

func report(command byte, payload ...byte) []byte {
	buf := make([]byte, reportLength)
	buf[0] = 0x11
	buf[1] = 0xff
	buf[2] = 0x04
	buf[3] = command
	copy(buf[4:], payload)
	return buf
}

func getSwitchOn() []byte {
	return report(0x1c, 0x01)
}

func getSwitchOff() []byte {
	return report(0x1c, 0x00)
}

func getLightState() []byte {
	return report(0x01)
}

func getSetBrightness(level int) []byte {
	level = min(max(level, 0), 100)

	value := minBrightness + (float64(level)/100)*(maxBrightness-minBrightness)

	return report(0x4c, 0x00, byte(value))
}

func getSetTemperature(temp int) []byte {
	temp = min(max(temp, minTemperature), maxTemperature)

	return report(0x9c, byte(temp>>8), byte(temp))
}

func (d *LitraDevice) write(data []byte) error {
	if _, err := d.dev.Write(data); err != nil {
		return fmt.Errorf("writing to device: %w", err)
	}

	response := make([]byte, reportLength)
	if _, err := d.dev.ReadWithTimeout(response, 5*time.Second); err != nil {
		return fmt.Errorf("reading device response: %w", err)
	}
	time.Sleep(30 * time.Millisecond)

	return nil
}

func (d *LitraDevice) TurnOn() error {
	return d.write(getSwitchOn())
}

func (d *LitraDevice) TurnOff() error {
	return d.write(getSwitchOff())
}

func (d *LitraDevice) IsOn() (bool, error) {
	if _, err := d.dev.Write(getLightState()); err != nil {
		return false, fmt.Errorf("writing to device: %w", err)
	}

	response := make([]byte, reportLength)
	n, err := d.dev.ReadWithTimeout(response, 5*time.Second)
	if err != nil {
		return false, fmt.Errorf("reading device response: %w", err)
	}
	if n < 5 {
		return false, fmt.Errorf("short device response: %d bytes", n)
	}

	return response[4] == 0x01, nil
}

func (d *LitraDevice) SetBrightness(level int) error {
	return d.write(getSetBrightness(level))
}

func (d *LitraDevice) SetTemperature(temp int) error {
	return d.write(getSetTemperature(temp))
}

func (d *LitraDevice) Close() error {
	err := d.dev.Close()
	if exitErr := hid.Exit(); err == nil {
		err = exitErr
	}
	return err
}
