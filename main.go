package main

import (
	"file-downloader/internal/app"
	"path/filepath"
)

var configPath = filepath.Join("config", "config.yaml")

func main() {
	app.Run(configPath)
}
