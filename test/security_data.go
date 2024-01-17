package test

import (
	"embed"
)

//go:embed resources
var res embed.FS

func SecurityResources() *embed.FS {
	return &res
}
