//go:build windows
// +build windows

package assets

import _ "embed"

var (
	//go:embed shutdown.bat.tmpl
	ShutdownBatTmpl []byte
	//go:embed startup.bat.tmpl
	StartupBatTmpl []byte
	//go:embed fangzhengshusong.ttf
	FontBytes []byte
)
