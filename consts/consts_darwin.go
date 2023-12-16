//go:build darwin
// +build darwin

package consts

import "path/filepath"

var (
	StartupBash  = filepath.Join(WorkDir, "startup.sh")
	ShutdownBash = filepath.Join(WorkDir, "shutdown.sh")
)

const (
	WorkDir = "/usr/local/yiService"
	MainExe = "yiService"
)
