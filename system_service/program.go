package system_service

import (
	"awesomeProject1/api"
	"awesomeProject1/consts"
	"bytes"
	"fmt"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"runtime"
)

type Program struct {
	service service.Service
	exit    chan struct{}
	cmd     *exec.Cmd
}

var stopCmd = "kill -9 `pidof yi_service`"

func init() {
	if runtime.GOOS == "windows" {
		stopCmd = "taskkill /f /t /im yi_service.exe"
	} else if runtime.GOOS == "darwin" {
		stopCmd = "launchctl stop  yi_service"
	}
}

func (p *Program) Start(s service.Service) error {
	fmt.Printf("starting %s service\n", runtime.GOOS)
	if service.Interactive() {
		fmt.Printf("Running in terminal \n")
	} else {
		fmt.Sprintf("Running under service manager. \n")
	}
	go p.run()
	return nil
}

func (p *Program) run() error {
	log.Infof("%s service running %s arguments %v\n", runtime.GOOS, service.Platform(), os.Args)
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("%v\n", r)
		}
	}()
	if err := api.ServiceHandler(); err != nil {
		log.Errorf("%v\n", err)
		return err
	}
	return nil
}

func (p *Program) Stop(s service.Service) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%v", err)
		}
	}()
	log.Infof("program stop service %s", consts.ServiceName)
	if p.exit != nil {
		close(p.exit)
	}
	if p != nil && p.cmd != nil && p.cmd.Process != nil {
		if err := p.cmd.Process.Kill(); err != nil {
			log.Errorf("%v\n", err)
		}
	}
	commandStop := exec.Command(stopCmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	commandStop.Stdout = &out
	commandStop.Stderr = &stderr
	if err := commandStop.Run(); err != nil {
		log.Errorf("%v\n", err)
	}
	if commandStop.Process != nil {
		commandStop.Process.Kill()
	}
	log.Infof("PROGRAM STOPPED: %s\n", out.String())
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
