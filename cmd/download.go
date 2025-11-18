package cmd

import (
	"errors"
	"fmt"
	"os"

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

		info, err := downloader.GetRemoteFileInfo(url)
		if err != nil {
			return err
		}

		fmt.Println("File info:")
		fmt.Println("Size:", info.Size)
		fmt.Println("Type:", info.ContentType)
		fmt.Println("Accept ranges:", info.AcceptRanges)

		chunks := downloader.SplitIntoChunks(info.Size, threads)

		fmt.Println("Chunks:")
		for _, ch := range chunks {
			fmt.Printf("\t#%d: %d - %d\n", ch.Index, ch.Start, ch.End)
		}

		var state *downloader.State

		if resume {
			st, err := downloader.LoadState(outputPath)
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if st.URL != url {
				return fmt.Errorf("state URL does not match")
			}
			if st.Size != info.Size {
				return fmt.Errorf("state size does not match file size on server")
			}

			fmt.Println("Resume mode: continuing unfinished download...")
			state = st
		} else {
			state = &downloader.State{
				URL:    url,
				Output: outputPath,
				Size:   info.Size,
				Chunks: make([]downloader.ChunkState, len(chunks)),
			}

			for i, ch := range chunks {
				state.Chunks[i] = downloader.ChunkState{
					Index: ch.Index,
					Start: ch.Start,
					End:   ch.End,
					Done:  false,
				}
			}

			if err := state.Save(); err != nil {
				return fmt.Errorf("failed to safve state: %w", err)
			}
		}

		toDownload := make([]downloader.Chunk, 0)

		for i, ch := range chunks {
			if !state.Chunks[i].Done {
				toDownload = append(toDownload, ch)
			}
		}

		fmt.Printf("Chunks to download: %d\n", len(toDownload))

		fmt.Println("Downloading started...")

		out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer out.Close()

		if err := out.Truncate(info.Size); err != nil {
			return err
		}

		fmt.Println("Concurrent chunk downloading started...")

		if err := downloader.DownloadChunksConcurrently(url, toDownload, out, state); err != nil {
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
