package cmd

import (
	"file-downloader/internal/logger"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	downloadUse          = "download"
	downloadShortMessage = "Download file by URL"
)

func newDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   downloadUse,
		Short: downloadShortMessage,
		RunE: func(cmd *cobra.Command, args []string) error {
			log := cmd.Context().Value("log").(*logger.Logger)

			url := args[0]

			log.Info("starting download", "url", url)

			fmt.Println(args[0]) // TODO: удалить + логика загрузки файла

			log.Info("download completed", "url", url)

			return nil
		},
	}
}
