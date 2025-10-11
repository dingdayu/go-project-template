package assets

import "embed"

//go:embed index.html
var IndexFS embed.FS

//go:embed all:dist
var DistFS embed.FS
