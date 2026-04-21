package web

import "embed"

//go:embed all:templates static
var Assets embed.FS
