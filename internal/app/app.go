package app

import (
	"file-downloader/cmd"
	"file-downloader/config"
	"file-downloader/internal/logger"
)

func Run(configPath string) {
	cfg := config.MustLoad(configPath)

	log := logger.New(cfg.Env)
	defer log.MustClose()

	cmd.Execute(log)
}
