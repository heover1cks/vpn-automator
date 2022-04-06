package main

import (
	"github.com/heover1cks/vpn-automator/client"
	"github.com/heover1cks/vpn-automator/config"
	"github.com/heover1cks/vpn-automator/options"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	opts := options.NewOptions()
	opts.AddFlags()
	err := opts.Parse()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}
	if opts.Quiet {
		log.SetLevel(log.FatalLevel)
	}
	cfg, err := config.LoadConfigFile(opts.ConfigFilePath)
	if err != nil {
		log.Fatal("fail to load config: ", err)
	}
	connectVPN(*cfg, opts.Alias)
}

func connectVPN(cfg config.Config, alias string) {
	cur, err := cfg.FindCurrentConfig(alias)
	if err != nil {
		log.Fatal("fail to find matching config: ", err)
	}
	switch cur.Client {
	case "bigip":
		var bc = &client.BigIPEdgeClient{
			Conf: *cur,
		}
		bc.BigIPEdgeClientSequence()
	}
}

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}
