package cmd

import (
	"errors"
	"fmt"

	"github.com/SALutHere/file_downloader/internal/downloader"
	"github.com/spf13/cobra"
)

var (
	outputPath string
	threads    int
	resume     bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <url>",
	Short: "Download file by specified URL",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("url is not specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		if outputPath == "" {
			return errors.New("you must specify the save path via --output")
		}

		fmt.Println("Downloading started...")

		if err := downloader.DownloadSimple(url, outputPath); err != nil {
			return err
		}

		fmt.Println("The file has been successfully downloaded.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(
		&outputPath,
		"output",
		"o",
		"",
		"Path to save file",
	)
	downloadCmd.Flags().IntVarP(
		&threads,
		"threads",
		"t",
		4,
		"Number of threads of downloading",
	)
	downloadCmd.Flags().BoolVar(
		&resume,
		"resume",
		false,
		"Allow file download to continue if there is a status",
	)

	downloadCmd.MarkFlagsMutuallyExclusive("resume")
	downloadCmd.Flags().SortFlags = true

	downloadCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: file_downloader download <url> [flags]")
		fmt.Println("\nFlags:")
		cmd.Flags().PrintDefaults()
	})
}
