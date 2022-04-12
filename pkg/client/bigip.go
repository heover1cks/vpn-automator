package client

import (
	"github.com/go-vgo/robotgo"
	"github.com/heover1cks/vpn-automator/config"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type BigIPEdgeClient struct {
	process string
	Conf    config.CurrentConfig
}

func (b *BigIPEdgeClient) DisconnectBigIPEdgeClientSequence() {
	b.killBigIPEdgeClient()
}

func (b *BigIPEdgeClient) ConnectBigIPEdgeClientSequence() {
	b.openBigIPEdgeClient()
	b.killBigIPEdgeClient()
	b.openBigIPEdgeClient()
	b.loginBigIPEdgeClient()
}

func (b *BigIPEdgeClient) openBigIPEdgeClient() {
	b.getPID()
	cmd := exec.Command("open", b.Conf.Location)
	err := cmd.Start()
	b.getPID()
	log.Info("client process started")
	if err != nil {
		log.Fatal("failed to start Big IP Edge Client: ", err)
	}
	robotgo.Sleep(1)
}

func isBigIPEdgeClientAlive() bool {
	pid, err := exec.Command("pgrep", "BIG-IP").Output()
	if err != nil {
		log.Fatal("failed to get status: ", err)
	}
	if string(pid) == "" {
		//log.Info("Big-IP Edge Client process not exists")
		return false
	} else {
		return true
	}
}

func (b *BigIPEdgeClient) getPID() {
	pid, err := exec.Command("pgrep", "BIG-IP").Output()
	if string(pid) == "" {
		log.Info("client process not exists")
		return
	}
	b.process = strings.Split(string(pid), "\n")[0]
	log.Info("pid: ", b.process)
	if err != nil {
		log.Fatal("failed to get PID: ", err)
	}
}

func (b *BigIPEdgeClient) killBigIPEdgeClient() {
	b.getPID()
	pid, err := strconv.Atoi(b.process)
	log.Info("pid: ", pid)
	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		log.Fatal("failed to kill Big IP Edge Client: ", err)
	}
	robotgo.Sleep(1)
}

func (b *BigIPEdgeClient) loginBigIPEdgeClient() {
	if err := robotgo.KeyTap("tab", "shift"); err != nil {
		log.Error(err)
	}
	robotgo.TypeStr(b.Conf.ID)
	robotgo.Sleep(1)
	if err := robotgo.KeyTap("tab"); err != nil {
		log.Error(err)
	}
	robotgo.Sleep(1)
	robotgo.TypeStr(b.Conf.PW)
	robotgo.MilliSleep(300)
	if err := robotgo.KeyTap("enter"); err != nil {
		log.Error(err)
	}
	robotgo.MilliSleep(300)
	robotgo.TypeStr(b.Conf.PW)
	if err := robotgo.KeyTap("enter"); err != nil {
		log.Error(err)
	}
}