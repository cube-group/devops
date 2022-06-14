package setting

import (
	"embed"
)

var _fs map[string]embed.FS

func initEmbed(fs map[string]embed.FS) {
	_fs = fs
}

func EmbedLocal() embed.FS {
	return _fs["local"]
}
