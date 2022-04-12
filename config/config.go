package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_BIGIP_LOCATION_MAC = "/Applications/BIG-IP Edge Client.app"
)

type ClientEnum int

const (
	WireGuard ClientEnum = iota
	BigIPEdge
)

func (s ClientEnum) String() string {
	switch s {
	case WireGuard:
		return "wg"
	case BigIPEdge:
		return "bigip"
	}
	return "unknown"
}

type Config struct {
	Accounts   []*DefaultAccountConfig `yaml:"accounts"`
	VPNClients []*VPNClientConfig      `yaml:"vpn_clients"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}

type DefaultAccountConfig struct {
	// Default Config
	Alias  string `yaml:"alias"`
	Client string `yaml:"client"`

	// Big-IP Edge Client Config
	ID string `yaml:"id,omitempty"`
	PW string `yaml:"pw,omitempty"`

	// WireGuard Config
	ServiceName string `yaml:"service_name,omitempty"`

	Status string
}

func (a *DefaultAccountConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain DefaultAccountConfig
	if err := unmarshal((*plain)(a)); err != nil {
		return err
	}
	return nil
}

func (c *Config) RetrieveAccountAliases() map[string]string {
	ret := map[string]string{}
	for _, acc := range c.Accounts {
		ret[acc.Alias] = acc.Status
	}
	return ret
}

type VPNClientConfig struct {
	Alias    string `yaml:"alias"`
	Location string `yaml:"location,omitempty"`
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
	// Default Config
	Alias  string
	Client string

	// Big-IP Edge Config
	ID       string
	PW       string
	Location string

	// WireGuard Config
	ServiceName string
}

func (c *Config) FindCurrentConfig(alias string) (*CurrentConfig, error) {
	cc := CurrentConfig{}
	if alias == "__NO_ALIAS_VALUE__" {
		return &cc, errors.New("VPN Alias required")
	}
	for _, acc := range c.Accounts {
		if acc.Alias == alias || acc.ID == alias {
			cc = CurrentConfig{
				Alias:  acc.Alias,
				Client: acc.Client,
			}
			err := cc.ParseClientName()
			if err != nil {
				return &cc, errors.New("no matching VPN client")
			}
			switch cc.Client {
			case "bigip":
				if acc.ID != "" && acc.PW != "" {
					cc.ID = acc.ID
					cc.PW = acc.PW
				} else {
					return &cc, errors.New("ID/PW for SSL VPN required")
				}
				cc.ID = acc.ID
				cc.PW = acc.PW
			case "wg":
				if acc.ServiceName != "" {
					cc.ServiceName = acc.ServiceName
				} else {
					return &cc, errors.New("wireguard service name required")
				}
			}
			break
		}
	}
	if cc == (CurrentConfig{}) {
		return &cc, errors.New("no matching alias")
	}
	for _, vpn := range c.VPNClients {
		if vpn.Alias == cc.Client {
			cc.Location = vpn.Location
		}
	}
	if cc.Location == "" && cc.Client == "bigip" {
		cc.Location = DEFAULT_BIGIP_LOCATION_MAC
		return &cc, errors.Wrapf(errors.New("location not specified"), "default is: %s", DEFAULT_BIGIP_LOCATION_MAC)
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
