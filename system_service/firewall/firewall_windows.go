//go:build windows
// +build windows

package firewall

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
	"xsky.com/downloader/consts"
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
