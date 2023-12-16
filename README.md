# 注册系统服务拉起进程（win/linux/darwin）
> 使用第三方包： go get "github.com/kardianos/service"
> 日志库：go get "github.com/sirupsen/logrus"
> - log "github.com/sirupsen/logrus"

## 1 初始化日志部分

## 2 导入三方Service库+编写serviceConfig（描述、可执行文件路径等）

### 2.1 consts
①consts/consts.go
②consts/consts_darwin.go
③consts/consts_windows.go

### 2.2 system_service

## 3 编写服务脚本（startup、shutdown脚本）
### 3.1 startup脚本
1. assets/startup.bat.tmpl
```bash
@echo off

@REM It is recommended to start here.
@REM Modify bash args below.
@REM "name" and "tags" must be alphanumeric sequence: [a-zA-Z-_.]
@REM Chinese chars are not supported here in bash file.
@REM
@REM Example:
@REM
@REM yiService.exe service start ^
@REM
@REM Startup from here
@REM
c:/yiService/yiService.exe service start


echo "yiService started!"
echo ""
pause 
```
2. assets/startup.sh.tmpl
```bash
#!/bin/bash

# It is recommended to start here.
# Modify bash args below.
# "name" and "tags" must be alphanumeric sequence: [a-zA-Z-_.]
# Chinese chars are not supported here in bash file.
#
# Example:
#
# yiService.exe service start \
#
# Startup from here
#
# launchctl start yiService
./yiService service start


echo yiService started!
ps aux |grep yiService
echo "" 
```


### 3.2 shutdown脚本
1. assets/shutdown.bat.tmpl
```bash
@echo off
@REM This command bash will stop the yiService windows service
@REM And then uninstall this service from operation system
@REM Configurations will be remained in directory c:/yiService/yiService on the disk.
@REM You can restart from those configurations in the near future.
@REM
c:/yiService/yiService.exe service stop
set "$process=yiService.exe"
for %%a in (%$process%) do tasklist /fi "imagename eq %%a"  | find "%%a" && taskkill /f /im %%a

echo shutdown yiService successfully!
pause 
```

2. assets/shutdown.sh.tmpl
```bash
#!/bin/bash

# This command bash will stop the yiService windows service
# And then uninstall this service from operation system
# Configurations will be remained in directory c:/yiService on the disk.
# You can restart from those configurations in the near future.
#

./yiService service stop
PID=$(ps -eaf | grep '/usr/local/yiService' | grep -v grep | awk '{print $2}')
if [[ "" !=  "$PID" ]]; then
  echo "killing $PID"
  kill -9 "$PID"
fi

echo shutdown yiService successfully!
```

## 4 编写extractFiles（执行exe文件之后，拷贝exe到服务工作目录）
```go
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
```

## 5 编写firewall部分+main函数部分（开放端口+由系统服务拉起进程）
### 5.1 firewall
1. system_service/firewall/firewall_darwin.go：
```go
//go:build darwin
// +build darwin

package firewall

import (
    "awesomeProject1/consts"
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "strconv"
    "strings"
)

func OpenPort() {
    log.Infof("darwin firewall checking\n")

    cmd0 := exec.Command("/usr/bin/nc", "-z", "127.0.0.1", fmt.Sprintf("%d", consts.ServicePort))
    log.Warnf("cmd0=%s\n", cmd0)
    stdout, err := cmd0.CombinedOutput()
    result := string(stdout)
    if err != nil {
       log.Infof("err=%v \n", err)
       return
    }
    if strings.Contains(result, "command not found") {
       fmt.Println("[warn]:", result)
       return
    }
    if strings.Contains(result, "not running") {
       fmt.Println("[warn]:", result)
       return
    }
    if strings.Contains(result, strconv.Itoa(consts.ServicePort)) {
       log.Warnf("%d already opened\n", consts.ServicePort)
       return
    }
    cmd := exec.Command("bash", "-c", fmt.Sprintf("firewall-cmd --zone=public --add-port=%d/tcp --permanent && firewall-cmd --reload", consts.ServicePort))
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
       log.Warnf("%s", stderr.String())
       log.Warnf("%v", err)
    }
    log.Warnf(out.String())
}
```
2. system_service/firewall/firewall_windows.go
```go
//go:build windows
// +build windows

package firewall

import (
    "bytes"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "runtime"
)

func OpenPort() {
    log.Infot("windows firewall checking")
    cmd := exec.Command("cmd", "/c", "netsh advfirewall firewall delete rule name=\"yiService\"")
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if runtime.GOOS == "windows" {
    }
    if err := cmd.Run(); err != nil {
       log.Errorf("%s", stderr.String())
       log.Errorf("%v", err)
    }
    cmd2 := exec.Command("cmd", "/c",
       fmt.Sprintf("netsh advfirewall firewall add rule name=\"yiService\" dir=in action=allow protocol=TCP localport=%d",
          consts.ServicePort,
       ))
    var out2 bytes.Buffer
    var stderr2 bytes.Buffer
    cmd2.Stdout = &out2
    cmd2.Stderr = &stderr2
    if runtime.GOOS == "windows" {
    }
    if err := cmd2.Run(); err != nil {
       log.Errorf("%s", stderr2.String())
       log.Errorf("%v", err)
    }
}
```
### 5.2 main函数
```go
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
```