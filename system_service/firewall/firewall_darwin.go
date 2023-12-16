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
