//go:build darwin
// +build darwin

package assets

import _ "embed"

var (
	//go:embed shutdown.sh.tmpl
	ShutdownBatTmpl []byte
	//go:embed startup.sh.tmpl
	StartupBatTmpl []byte
	//go:embed SourceHanSansCN-Normal.ttf
	FontBytes []byte
)
