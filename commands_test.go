package main

import (
	"testing"
)

func TestIntArg(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    int
		wantErr bool
	}{
		{"valid", []string{"50"}, 50, false},
		{"lower bound", []string{"0"}, 0, false},
		{"upper bound", []string{"100"}, 100, false},
		{"below range", []string{"-1"}, 0, true},
		{"above range", []string{"101"}, 0, true},
		{"not a number", []string{"abc"}, 0, true},
		{"no args", []string{}, 0, true},
		{"too many args", []string{"50", "60"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := intArg(tt.args, "brightness", 0, 100)
			if (err != nil) != tt.wantErr {
				t.Fatalf("intArg(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("intArg(%v) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}

func TestCommandNames(t *testing.T) {
	want := []string{"on", "off", "toggle", "brightness", "temperature", "brightness-up", "brightness-down", "temperature-up", "temperature-down", "profile"}

	if len(commands) != len(want) {
		t.Fatalf("len(commands) = %d, want %d", len(commands), len(want))
	}

	seen := make(map[string]bool)
	for i, cmd := range commands {
		if cmd.name != want[i] {
			t.Errorf("commands[%d].name = %q, want %q", i, cmd.name, want[i])
		}
		if seen[cmd.name] {
			t.Errorf("duplicate command name %q", cmd.name)
		}
		seen[cmd.name] = true
		if cmd.run == nil {
			t.Errorf("command %q has no run function", cmd.name)
		}
	}
}

func TestOptionalStepArg(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		defaultVal int
		want       int
		wantErr    bool
	}{
		{"no args uses default", []string{}, 10, 10, false},
		{"explicit step", []string{"5"}, 10, 5, false},
		{"step of 1", []string{"1"}, 10, 1, false},
		{"zero step", []string{"0"}, 10, 0, true},
		{"negative step", []string{"-1"}, 10, 0, true},
		{"not a number", []string{"abc"}, 10, 0, true},
		{"too many args", []string{"5", "10"}, 10, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := optionalStepArg(tt.args, "test", tt.defaultVal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("optionalStepArg(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("optionalStepArg(%v) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	if err := run(nil); err == nil {
		t.Error("run(nil) = nil, want error for missing command")
	}

	if err := run([]string{"sparkle"}); err == nil {
		t.Error(`run(["sparkle"]) = nil, want error for unknown command`)
	}
}
