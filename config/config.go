package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Config struct {
	Accounts   []*AccountConfig   `yaml:"accounts"`
	VPNClients []*VPNClientConfig `yaml:"vpn_clients"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

type AccountConfig struct {
	Alias  string `yaml:"alias"`
	ID     string `yaml:"id"`
	PW     string `yaml:"pw"`
	Client string `yaml:"client"`
}

func (a *AccountConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain AccountConfig
	if err := unmarshal((*plain)(a)); err != nil {
		return err
	}
	return nil
}

type VPNClientConfig struct {
	Alias    string `yaml:"alias"`
	Location string `yaml:"location"`
}

func (v *VPNClientConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain VPNClientConfig
	if err := unmarshal((*plain)(v)); err != nil {
		return err
	}
	return nil
}

func LoadConfigFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := loadConfig(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	cfg.SetDirectory(filepath.Dir(filename))
	return cfg, nil
}

func loadConfig(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) SetDirectory(dir string) {}

type CurrentConfig struct {
	Alias    string
	ID       string
	PW       string
	Client   string
	Location string
}

func (c *Config) FindCurrentConfig(alias string) (*CurrentConfig, error) {
	cc := CurrentConfig{}
	if alias == "__NO_ALIAS_VALUE__" {
		return &cc, errors.New("VPN Alias required")
	}
	for _, acc := range c.Accounts {
		if acc.Alias == alias {
			cc = CurrentConfig{
				Alias:  acc.Alias,
				ID:     acc.ID,
				PW:     acc.PW,
				Client: acc.Client,
			}
		}
	}
	cc.ParseClientName()

	if cc == (CurrentConfig{}) {
		return &cc, errors.New("no matching alias")
	}
	for _, vpn := range c.VPNClients {
		if vpn.Alias == cc.Client {
			cc.Location = vpn.Location
		}
	}
	if cc.Location == "" {
		return &cc, errors.New("no matching VPN client")
	}
	return &cc, nil
}

func (cc *CurrentConfig) ParseClientName() error {
	temp := strings.ToLower(cc.Client)
	temp = strings.ReplaceAll(temp, " ", "")
	temp = strings.ReplaceAll(temp, "-", "")
	temp = strings.ReplaceAll(temp, "_", "")

	if strings.Contains(temp, "bigip") {
		cc.Client = "bigip"
		return nil
	} else if strings.Contains(temp, "wireguard") || temp == "wg" {
		cc.Client = "wg"
		return nil
	} else {
		cc.Client = ""
		return errors.New("not a valid vpn client")
	}
}
