package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"litra/driver"
)

type Profile struct {
	Brightness  int `json:"brightness"`
	Temperature int `json:"temperature"`
}

type Config struct {
	Profiles map[string]Profile `json:"profiles"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".litra"), nil
}

func loadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	return &cfg, nil
}

func (c *Config) saveToPath(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	return loadConfigFromPath(path)
}

func (c *Config) save() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	return c.saveToPath(path)
}

func profileSave(ld *driver.LitraDevice, name string) error {
	state, err := ld.GetState()
	if err != nil {
		return fmt.Errorf("reading device state: %w", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	cfg.Profiles[name] = Profile{
		Brightness:  state.Brightness,
		Temperature: state.Temperature,
	}

	if err := cfg.save(); err != nil {
		return err
	}

	fmt.Printf("Profile %q saved (brightness: %d%%, temperature: %dK)\n", name, state.Brightness, state.Temperature)
	return nil
}

func profileLoad(ld *driver.LitraDevice, name string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	p, ok := cfg.Profiles[name]
	if !ok {
		return fmt.Errorf("profile %q not found", name)
	}

	if err := ld.SetBrightness(p.Brightness); err != nil {
		return fmt.Errorf("setting brightness: %w", err)
	}
	if err := ld.SetTemperature(p.Temperature); err != nil {
		return fmt.Errorf("setting temperature: %w", err)
	}

	fmt.Printf("Profile %q loaded (brightness: %d%%, temperature: %dK)\n", name, p.Brightness, p.Temperature)
	return nil
}

func profileList() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles saved")
		return nil
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Profiles:")
	for _, name := range names {
		p := cfg.Profiles[name]
		fmt.Printf("  %-20s brightness: %d%%, temperature: %dK\n", name, p.Brightness, p.Temperature)
	}
	return nil
}
