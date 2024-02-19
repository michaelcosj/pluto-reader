package assets

import (
	"embed"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

func Mount() http.Handler {
    return http.FileServer(http.FS(dist))
}
