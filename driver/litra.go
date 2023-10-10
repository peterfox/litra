package driver

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	vendor  = 0x046d
	product = 0xc900
)

type LitraDevice struct {
	dev *hid.Device
}

func New() (*LitraDevice, error) {
	d := &LitraDevice{}

	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	dev, err := hid.OpenFirst(vendor, product)
	d.dev = dev

	return d, err
}

func getSwitchOn() []byte {
	return []byte{0x11, 0xff, 0x04, 0x1c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func getSwitchOff() []byte {
	return []byte{0x11, 0xff, 0x04, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func getLightState() []byte {
	return []byte{0x11, 0xff, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func getSetBrightness(level int) []byte {
	minBrightness := float64(0x14)
	maxBrightness := float64(0xfa)

	if level < 0 {
		level = 0
	}
	if level > 100 {
		level = 100
	}

	value := minBrightness + ((float64(level) / 100) * (maxBrightness - minBrightness))
	adjusted_level := byte(math.Floor(float64(value)))

	return []byte{0x11, 0xff, 0x04, 0x4c, 0x00, adjusted_level, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func getSetTemperature(temp int16) []byte {
	if temp < 2700 {
		temp = 2700
	}
	if temp > 6500 {
		temp = 6500
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, temp)
	if err != nil {
		return nil
	}
	byte0, _ := buf.ReadByte()
	byte1, _ := buf.ReadByte()

	return []byte{0x11, 0xff, 0x04, 0x9c, byte0, byte1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func (d *LitraDevice) TurnOn() error {
	dummy := []byte{0x00}

	_, err := d.dev.Write(getSwitchOn())
	_, err = d.dev.Read(dummy)
	time.Sleep(30 * time.Millisecond)

	return err
}

func (d *LitraDevice) IsOn() (bool, error) {
	dummy := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	_, err := d.dev.Write(getLightState())
	_, err = d.dev.ReadWithTimeout(dummy, time.Second*5)

	if err != nil {
		panic(err)
	}

	return dummy[4] == 0x01, err
}

func (d *LitraDevice) TurnOff() error {
	dummy := []byte{0x00}

	_, err := d.dev.Write(getSwitchOff())
	_, err = d.dev.Read(dummy)
	time.Sleep(30 * time.Millisecond)

	return err
}

func (d *LitraDevice) SetBrightness(level int) error {
	dummy := []byte{0x00}

	_, err := d.dev.Write(getSetBrightness(level))
	_, err = d.dev.Read(dummy)
	time.Sleep(30 * time.Millisecond)

	return err
}

func (d *LitraDevice) SetTemperature(temp int16) error {
	dummy := []byte{0x00}

	_, err := d.dev.Write(getSetTemperature(temp))
	_, err = d.dev.Read(dummy)
	time.Sleep(30 * time.Millisecond)

	return err
}

func (d *LitraDevice) Close() error {
	return d.dev.Close()
}
