package data

import (
	"embed"
	_ "embed"
)

//go:embed *.yaml
var DataFS embed.FS
