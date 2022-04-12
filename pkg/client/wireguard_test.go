package client

import (
	"github.com/heover1cks/vpn-automator/config"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var expectedConfig = &config.CurrentConfig{
	Alias:  "kcloud",
	Client: "wg",
}

func TestWireGuardClient_WireGuardClientSequence(t *testing.T) {
	wg := WireGuardClient{
		Conf: *expectedConfig,
	}
	wg.getStatus()
	assert.Contains(t, strings.ToLower(wg.Status), "connected")
}
