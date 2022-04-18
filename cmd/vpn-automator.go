package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/heover1cks/vpn-automator/config"
	"github.com/heover1cks/vpn-automator/pkg/client"
	"github.com/heover1cks/vpn-automator/pkg/options"
	"github.com/heover1cks/vpn-automator/pkg/tui"
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

	if opts.CLI {
		cfg, err := config.LoadConfigFile(opts.ConfigFilePath)
		if err != nil {
			log.Fatal("fail to load config: ", err)
		}
		client.ConnectVPN(*cfg, opts.Alias, opts.Disconnect)

	} else {
		cfg, err := config.LoadConfigFile(os.Getenv("VPN_AUTOMATOR_CONFIG_PATH"))
		if err != nil {
			log.Fatal("fail to load config: ", err)
		}
		p := tea.NewProgram(tui.InitialModel(*cfg))
		if err := p.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(os.Getenv("VPN_AUTOMATOR_REPORT_CALLER") == "true")
}
