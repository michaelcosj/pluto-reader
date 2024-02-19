package assets

import (
	"embed"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

func Mount(prefix string) http.Handler {
    return http.StripPrefix(prefix, http.FileServer(http.FS(dist)))
}
