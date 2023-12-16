package consts

import "path/filepath"

var (
	StartupBash  = filepath.Join(WorkDir, "startup.bat")
	ShutdownBash = filepath.Join(WorkDir, "shutdown.bat")
)

const (
	WorkDir = "c:/yiService"
	MainExe = "yiService.exe"
)
