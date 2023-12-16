package main

import (
	"awesomeProject1/consts"
	"awesomeProject1/system_service"
	"awesomeProject1/system_service/firewall"
	"awesomeProject1/util"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	//os.Setenv("dev", "true")
	util.InitLog(consts.LogPath)
}

func main() {
	if len(os.Getenv("dev")) != 0 {
		system_service.StopSVC()
		firewall.OpenPort()
		system_service.StartSVC()
		system_service.RunSVC()
	}
	log.Errorf("os.Args=%v len=%v \n", os.Args, len(os.Args))
	if len(os.Args) == 1 {
		//stop svc if exist
		system_service.StopSVC()
		log.Errorf("install %v \n", consts.ServiceName)
		if err := os.MkdirAll(consts.WorkDir, os.ModePerm); err != nil {
			log.Errorf("%v\n", err)
		}
		firewall.OpenPort()
		system_service.ExtractFiles()
		pwd, err := os.Getwd()
		if err != nil {
			log.Errorf("%v\n", err)
		}
		log.Infof("install svc, working directory %s", pwd)
		system_service.StartSVC()
		log.Infof("yiService installed!")
		return
	}
	os.Chdir(consts.WorkDir)
	log.Errorf("service %s \n", os.Args[2])
	switch os.Args[2] {
	case "start":
		system_service.StartSVC()
		return
	case "stop":
		system_service.StopSVC()
		return
	default:
		system_service.RunSVC()
		log.Info("running yiService")
	}
}
