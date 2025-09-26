/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
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
		whisperService := service.NewScriptServiceImpl()

		var targetWidth = 1080
		var targetHeight = 1920
		var fadeDuration = 5

		fmt.Println("Verbose: ", verbose)

		var clipQueue []*model.ClipDTO

		clipQueue = append(clipQueue, model.NewClipDTO(
			audioPath,
			videoPath,
			nil,
			nil,
			nil,
		))

		if !helper.IsValidClipQueue(clipQueue) {
			return errors.New("found no files to analyze")
		}

		if err := helper.CreateDirectoryIfNotExists(outputDir); err != nil {
			return errors.New(fmt.Sprintf("couldn't create output directory: %s", outputDir))
		}

		for _, clip := range clipQueue {
			if !clip.IsValidAudioInputPath() {
				return errors.New(fmt.Sprintf("%s is not a valid input path", clip.AudioInputPath))
			}

			fmt.Println("Starting SRT generation...")
			if err := helper.GenerateSRTCaptions(whisperService, outputDir, clip, whisperModel, startTime, endTime, verbose); err != nil {
				return err
			}

			fmt.Println("Starting burn...")
			if err := helper.BurnCaptions(whisperService, outputDir, clip, &targetWidth, &targetHeight, startTime, endTime, verbose); err != nil {
				return err
			}

			fmt.Println("Starting trim and fade...")
			if err := helper.TrimAndFade(whisperService, outputDir, clip, "00:00", "00:20", &fadeDuration, verbose); err != nil {
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
}
