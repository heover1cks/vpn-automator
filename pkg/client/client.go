package client

import (
	"github.com/heover1cks/vpn-automator/config"
	log "github.com/sirupsen/logrus"
)

const (
	// TODO: Change to icon
	Connected        = "CONN"
	Disconnected     = "DISC"
	ConnectedUnknown = "UNKN"
)

func ConnectVPN(cfg config.Config, alias string, disconnect bool) {
	cur, err := cfg.FindCurrentConfig(alias)
	if err != nil {
		log.Fatal("fail to find matching config: ", err)
	}
	switch cur.Client {
	case "bigip":
		var bc = &BigIPEdgeClient{
			Conf: *cur,
		}
		if !disconnect {
			bc.ConnectBigIPEdgeClientSequence()
		} else {
			bc.DisconnectBigIPEdgeClientSequence()
		}
	case "wg":
		var wg = &WireGuardClient{
			Conf: *cur,
		}
		if !disconnect {
			wg.ConnectWireGuardClientSequence()
		} else {
			wg.DisconnectWireGuardClientSequence()
		}
	}
}

func GetConnectionStatus(conf config.Config) config.Config {
	for _, acc := range conf.Accounts {
		switch acc.Client {
		case "bigip":
			if isBigIPEdgeClientAlive() {
				acc.Status = ConnectedUnknown
			} else {
				acc.Status = Disconnected
			}
		case "wg":
			if IsWireguardAlive(acc.ServiceName) {
				acc.Status = Connected

			} else {
				acc.Status = Disconnected
			}
		}
	}
	return conf
}
