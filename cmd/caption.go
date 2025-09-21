/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/helper"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/service"
	"github.com/spf13/cobra"
)

// captionCmd represents the caption command
var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: "Generates captions for an audio file.",
	Long:  `Generates captions for an audio file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		whisperService := service.NewWhisperServiceImpl()

		fmt.Println("Verbose: ", verbose)

		var clipQueue []*model.ClipDTO

		if isDirectory {
			if !helper.IsDirectory(path) {
				return errors.New(fmt.Sprintf("%s is not a directory", path))
			}

			files, err := helper.GetFilesInDirectory(path)
			if err != nil {
				return errors.New(fmt.Sprintf("couldn't get files in %s", path))
			}

			clipQueue = helper.InputPathArrayToClips(files)
		} else {
			clipQueue = append(clipQueue, helper.ClipFromInputPath(path))
		}

		if !helper.IsValidClipQueue(clipQueue) {
			return errors.New("found no files to analyze")
		}

		if err := helper.CreateDirectoryIfNotExists(outputDir); err != nil {
			return errors.New(fmt.Sprintf("couldn't create output directory: %s", outputDir))
		}

		for _, clip := range clipQueue {
			if !clip.IsValidInputPath() {
				return errors.New(fmt.Sprintf("%s is not a valid input path", *clip.InputPath))
			}

			if err := helper.GenerateSRTCaptions(whisperService, outputDir, clip, verbose); err != nil {
				return err
			}
		}

		for _, clip := range clipQueue {
			if err := clip.PrintTable(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(captionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// captionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// captionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
