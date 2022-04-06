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
	flags          *pflag.FlagSet
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
	o.flags.BoolVarP(&o.Quiet, "quiet", "q", false, "Running quietly (default: false, print info logs)")
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