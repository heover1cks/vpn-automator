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
	if opts.LogFormat != "default" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	cfg, err := config.LoadConfigFile(opts.ConfigFilePath)
	if err != nil {
		log.Fatal("fail to load config: ", err)
	}
	connectVPN(*cfg, opts.Alias, opts.Disconnect)
}

func connectVPN(cfg config.Config, alias string, disconnect bool) {
	cur, err := cfg.FindCurrentConfig(alias)
	if err != nil {
		log.Fatal("fail to find matching config: ", err)
	}
	switch cur.Client {
	case "bigip":
		var bc = &client.BigIPEdgeClient{
			Conf: *cur,
		}
		if !disconnect {
			bc.ConnectBigIPEdgeClientSequence()
		} else {
			bc.DisconnectBigIPEdgeClientSequence()
		}
	case "wg":
		var wg = &client.WireGuardClient{
			Conf: *cur,
		}
		if !disconnect {
			wg.ConnectWireGuardClientSequence()
		} else {
			wg.DisconnectWireGuardClientSequence()
		}
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}
