package system_service

import (
	"awesomeProject1/assets"
	"awesomeProject1/consts"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
)

func ExtractFiles() {
	CopyMainExe()
	//copy startup\shutdown
	if err := os.WriteFile(consts.StartupBash, assets.StartupBatTmpl, os.ModePerm); err != nil {
		log.Errorf("create startup bash failed %v", err)
	}
	if err := os.WriteFile(consts.ShutdownBash, assets.ShutdownBatTmpl, os.ModePerm); err != nil {
		log.Errorf("create shutdown bash failed %v", err)
	}
}

func CopyMainExe() {
	executable, err := os.Executable()
	log.Infof("install %s to %s", executable, consts.MainExe)
	if err != nil {
		log.Errorf("%v", err)
	}
	sFile, err := os.Open(executable)
	if err != nil {
		log.Errorf("%v", err)
	}
	defer sFile.Close()
	exePath := fmt.Sprintf("%s/%s", consts.WorkDir, consts.MainExe)
	if runtime.GOOS == "windows" {
		exePath = fmt.Sprintf("%s\\%s", consts.WorkDir, consts.MainExe)
	}
	_, err = os.Stat(exePath)
	if err == nil {
		//overwrite
		if err := os.RemoveAll(exePath); err != nil {
			log.Errorf("%v", err)
		}
	}
	eFile, err := os.Create(exePath)
	if err != nil {
		log.Errorf("%v", err)
	}
	defer eFile.Close()
	if _, err = io.Copy(eFile, sFile); err != nil {
		log.Errorf("%v", err)
	}
	if err = eFile.Sync(); err != nil {
		log.Errorf("%v", err)
	}
	if err = os.Chdir(consts.WorkDir); err != nil {
		log.Errorf("%v\n", err)
	}
	if err = os.Chmod(consts.MainExe, os.FileMode(0777)); err != nil {
		log.Errorf("%v", err)
	}
}
