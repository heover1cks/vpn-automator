package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type Options struct {
	Help           bool
	ConfigFilePath string
	Alias          string
	Quiet          bool
	Disconnect     bool
	LogFormat      string

	CLI   bool
	flags *pflag.FlagSet
}

func (o *Options) AddFlags() {
	o.flags = pflag.NewFlagSet("", pflag.ExitOnError)
	o.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		o.flags.PrintDefaults()
	}
	o.flags.BoolVarP(&o.Help, "help", "h", false, "Print Help text")
	o.flags.StringVarP(&o.Alias, "alias", "a", "__NO_ALIAS_VALUE__", "Select VPN alias to connect(required)")
	o.flags.StringVarP(&o.ConfigFilePath, "config", "c", "./config.yml", "Config file path (default: ./config.yml)")
	o.flags.StringVarP(&o.LogFormat, "format", "f", "default", "Log Format (default: text, available: json)")
	o.flags.BoolVarP(&o.Quiet, "quiet", "q", false, "Running quietly (default: false, print info logs)")
	o.flags.BoolVarP(&o.Disconnect, "disconnect", "d", false, "Disconnect (default: false)")

	o.flags.BoolVar(&o.CLI, "cli", false, "Run as CLI Mode (default: false, run as TUI mode)")

}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Parse() error {
	err := o.flags.Parse(os.Args)
	return err
}

func (o *Options) Usage() {
	o.flags.Usage()
}
