//go:build windows
// +build windows

package downloader

import (
	"bytes"
	"fmt"
	"github.com/aobco/log"
	"os/exec"
	"runtime"
	"xsky.com/downloader/consts"
)

func Firewall() {
	log.Infof("windows firewall checking")
	cmd := exec.Command("cmd", "/c", "netsh advfirewall firewall delete rule name=\"xDownload\"")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if runtime.GOOS == "windows" {
	}
	if err := cmd.Run(); err != nil {
		log.Warnf("%s", stderr.String())
		log.Warnf("%v", err)
	}
	cmd2 := exec.Command("cmd", "/c",
		fmt.Sprintf("netsh advfirewall firewall add rule name=\"xDownload\" dir=in action=allow protocol=TCP localport=%d",
			consts.XDownPort,
		))
	var out2 bytes.Buffer
	var stderr2 bytes.Buffer
	cmd2.Stdout = &out2
	cmd2.Stderr = &stderr2
	if runtime.GOOS == "windows" {
	}
	if err := cmd2.Run(); err != nil {
		log.Warnf("%s", stderr2.String())
		log.Warnf("%v", err)
	}
}
