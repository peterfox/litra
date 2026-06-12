package driver

import (
	"bytes"
	"testing"
)

func TestReport(t *testing.T) {
	got := report(0x1c, 0x01)

	if len(got) != reportLength {
		t.Fatalf("report length = %d, want %d", len(got), reportLength)
	}

	want := []byte{0x11, 0xff, 0x04, 0x1c, 0x01}
	if !bytes.Equal(got[:5], want) {
		t.Errorf("report header = %#v, want %#v", got[:5], want)
	}

	for i, b := range got[5:] {
		if b != 0x00 {
			t.Errorf("report[%d] = %#x, want zero padding", i+5, b)
		}
	}
}

func TestSwitchReports(t *testing.T) {
	if got := getSwitchOn(); got[3] != 0x1c || got[4] != 0x01 {
		t.Errorf("getSwitchOn() = %#v, want command 0x1c with payload 0x01", got[:5])
	}
	if got := getSwitchOff(); got[3] != 0x1c || got[4] != 0x00 {
		t.Errorf("getSwitchOff() = %#v, want command 0x1c with payload 0x00", got[:5])
	}
	if got := getLightState(); got[3] != 0x01 || got[4] != 0x00 {
		t.Errorf("getLightState() = %#v, want command 0x01 with no payload", got[:5])
	}
}

func TestGetSetBrightness(t *testing.T) {
	tests := []struct {
		level int
		want  byte
	}{
		{0, minBrightness},
		{100, maxBrightness},
		{50, 0x87},
		{-10, minBrightness},
		{200, maxBrightness},
	}

	for _, tt := range tests {
		got := getSetBrightness(tt.level)
		if got[3] != 0x4c {
			t.Errorf("getSetBrightness(%d) command = %#x, want 0x4c", tt.level, got[3])
		}
		if got[5] != tt.want {
			t.Errorf("getSetBrightness(%d) value = %#x, want %#x", tt.level, got[5], tt.want)
		}
	}
}

func TestGetSetTemperature(t *testing.T) {
	tests := []struct {
		temp int
		want uint16
	}{
		{2700, 2700},
		{6500, 6500},
		{4000, 4000},
		{0, minTemperature},
		{10000, maxTemperature},
	}

	for _, tt := range tests {
		got := getSetTemperature(tt.temp)
		if got[3] != 0x9c {
			t.Errorf("getSetTemperature(%d) command = %#x, want 0x9c", tt.temp, got[3])
		}
		if value := uint16(got[4])<<8 | uint16(got[5]); value != tt.want {
			t.Errorf("getSetTemperature(%d) value = %d, want %d", tt.temp, value, tt.want)
		}
	}
}
