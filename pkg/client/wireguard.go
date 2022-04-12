package client

import (
	"github.com/heover1cks/vpn-automator/config"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type WireGuardClient struct {
	Conf   config.CurrentConfig
	Status string
}

func (w *WireGuardClient) ConnectWireGuardClientSequence() {
	w.getStatus()
	if w.Status == "connected" {
		log.Info("already connected to vpn: ", w.Conf.ServiceName)
		return
	}
	w.connect()
}
func (w *WireGuardClient) DisconnectWireGuardClientSequence() {
	w.getStatus()
	if w.Status == "disconnected" {
		log.Info("already disconnected to vpn: ", w.Conf.ServiceName)
		return
	}
	w.disconnect()
}

func IsWireguardAlive(service string) bool {
	println(service)
	status, err := exec.Command("scutil", "--nc", "status", service).Output()
	if err != nil {
		log.Fatal("failed to get wireguard-network status: ", err)
	}
	return strings.ToLower(strings.Split(string(status), "\n")[0]) == "connected"
}

func (w *WireGuardClient) getStatus() {
	status, err := exec.Command("scutil", "--nc", "status", w.Conf.ServiceName).Output()
	if err != nil {
		log.Fatal("failed to get wireguard-network status: ", err)
	}
	w.Status = strings.ToLower(strings.Split(string(status), "\n")[0])
}

func (w *WireGuardClient) connect() {
	cmd := exec.Command("scutil", "--nc", "start", w.Conf.ServiceName)
	err := cmd.Start()
	if err != nil {
		log.Fatal("failed to start wireguard-network: ", w.Conf.ServiceName, err)
	}
	w.getStatus()
	log.Info("Connected to WireGuard: ", w.Conf.ServiceName)
}

func (w *WireGuardClient) disconnect() {
	cmd := exec.Command("scutil", "--nc", "stop", w.Conf.ServiceName)
	err := cmd.Start()
	if err != nil {
		log.Fatal("failed to stop wireguard-network: ", w.Conf.ServiceName, err)
	}
	w.getStatus()
	log.Info("Disconnected to WireGuard: ", w.Conf.ServiceName)
}
