package consts

import "path"

const (
	ServiceVersion     = "v1.0"
	ServiceName        = "yiService"
	ServicePort        = 9999
	ServiceDisplayName = "yi.service"
	ServiceDescription = "my test service"
)

var LogPath string

func init() {
	LogPath = path.Join(WorkDir, "yiService.log")
}
