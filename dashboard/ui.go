package dashboard

import (
	"embed"
	"path"
)

//go:embed dist/*
var EmbedFS embed.FS

func Open(p string) ([]byte, error) {
	return EmbedFS.ReadFile(path.Join("dist", p))
}
