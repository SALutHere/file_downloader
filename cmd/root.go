package cmd

import (
	"context"
	"file-downloader/internal/logger"
	"os"

	"github.com/spf13/cobra"
)

const (
	rootShortMessage = "Trusted file downloader"
	rootLongMessage  = "CLI-tool for fast and stable file loading with " +
		"support of multithreading and pausing resumption of downloading."
)

func Execute(log *logger.Logger) {
	rootCmd := newRootCmd(log)
	if err := rootCmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func newRootCmd(log *logger.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Short: rootShortMessage,
		Long:  rootLongMessage,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(cmd.Context(), "log", log)
			cmd.SetContext(ctx)
			log.Info("CLI initialized")
		},
	}

	rootCmd.AddCommand(newDownloadCmd())

	return rootCmd
}
