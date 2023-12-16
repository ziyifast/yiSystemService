package system_service

import (
	"awesomeProject1/consts"
	"fmt"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var (
	logger service.Logger
	svc    service.Service
)

func init() {
	svcConfig := &service.Config{
		Name:             consts.ServiceName,
		WorkingDirectory: consts.WorkDir,
		DisplayName:      consts.ServiceDisplayName,
		Description:      consts.ServiceDescription,
		Arguments:        []string{"service", "run"}, //服务注册成功之后，由服务去执行yiService.exe service run【运行服务】
		Executable:       fmt.Sprintf("%s/%s", consts.WorkDir, consts.ServiceName),
		Option:           service.KeyValue{},
	}
	if runtime.GOOS == "windows" {
		svcConfig.Executable = fmt.Sprintf("%s\\%s.exe", consts.WorkDir, consts.ServiceName)
	}
	var err error
	program := &Program{}
	svc, err = service.New(program, svcConfig)
	if err != nil {
		log.Errorf("create service fail %v\n", err)
		return
	}
	errChan := make(chan error, 5)
	logger, err = svc.Logger(errChan)
	if err != nil {
		log.Errorf("%v\n", err)
		return
	}
	if err != nil {
		log.Errorf("%v\n", err)
		return
	}
	go func() {
		log.Info("watching err chan....")
		for {
			err := <-errChan
			if err != nil {
				log.Fatalf("service err %v", err)
			}
		}
	}()
}

func StartSVC() {
	log.Infof("StartSVC...")
	serviceControl("install")
	serviceControl("start")
}

func StopSVC() {
	log.Infof("try to stop service, if already exists.")
	serviceControl("stop")
}

func RunSVC() {
	fmt.Sprintf("%s service running \n", runtime.GOOS)
	if err := svc.Run(); err != nil {
		fmt.Sprintf("%s service running fail %v \n", runtime.GOOS, err)
		os.Exit(1)
	}
}

func serviceControl(action string) {
	log.Infof("%s service %s \n", runtime.GOOS, action)
	if err := service.Control(svc, action); err != nil {
		log.Infof("%s service: %v \n", action, err)
	}
}
