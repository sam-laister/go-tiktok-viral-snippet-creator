/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/helper"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/service"
	"github.com/spf13/cobra"
)

var captionsOptions = model.NewCaptionOptions()

var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: "Generates on demand captions for an audio file/directory",
	Long:  `Captions doesn't use an external database and instead acts as a purely I/O caption generator. Files are generated using timestamp and not metadata.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		whisperService := service.NewScriptServiceImpl()

		fmt.Println("Verbose: ", captionsOptions.Verbose)

		var clipQueue []*model.ClipDTO

		if captionsOptions.IsDirectory {
			if !helper.IsDirectory(captionsOptions.AudioPath) || !helper.IsDirectory(captionsOptions.VideoPath) {
				return errors.New(fmt.Sprintf(
					"either %s or %s is not a directory",
					captionsOptions.AudioPath,
					captionsOptions.VideoPath,
				))
			}

			audios, err := helper.GetFilesInDirectory(captionsOptions.AudioPath)
			if err != nil {
				return err
			}

			videos, err := helper.GetFilesInDirectory(captionsOptions.VideoPath)
			if err != nil {
				return err
			}

			// Random video with each audio
			for _, audio := range audios {
				clipQueue = append(clipQueue, model.NewClipDTO(
					audio,
					videos[rand.Intn(len(videos))],
					nil,
					nil,
					nil,
					nil,
					nil,
				))
			}

		} else {
			clipQueue = append(clipQueue, model.NewClipDTO(
				captionsOptions.AudioPath,
				captionsOptions.VideoPath,
				nil,
				nil,
				nil,
				nil,
				nil,
			))
		}

		if !helper.IsValidClipQueue(clipQueue) {
			return errors.New("found no files to analyze")
		}

		if err := helper.CreateDirectoryIfNotExists(captionsOptions.OutputDir); err != nil {
			return errors.New(fmt.Sprintf(
				"couldn't create output directory: %s",
				captionsOptions.OutputDir,
			))
		}

		for index, clip := range clipQueue {
			fmt.Println(fmt.Sprintf("Processing batch %d/%d", index, len(clipQueue)))

			if !clip.IsValidAudioInputPath() {
				return errors.New(fmt.Sprintf("%s is not a valid input path", clip.AudioInputPath))
			}

			fmt.Println("Starting SRT generation...")
			if err := whisperService.RunGenerateSRTCaptionsOnClip(
				captionsOptions.OutputDir,
				clip,
				captionsOptions.WhisperModel,
				captionsOptions.StartTime,
				captionsOptions.EndTime,
				captionsOptions.Verbose,
			); err != nil {
				return err
			}

			if !captionsOptions.NoInteract {
				if err := clip.PrintTable(); err != nil {
					return err
				}

				if _, err := helper.WaitForOptionalEdits(); err != nil {
					return err
				}
			}

			fmt.Println("Starting burn...")
			if err := whisperService.RunBurnCaptionsOnClip(
				captionsOptions.OutputDir,
				clip,
				&captionsOptions.Width,
				&captionsOptions.Height,
				captionsOptions.StartTime,
				captionsOptions.EndTime,
				captionsOptions.Verbose,
			); err != nil {
				return err
			}

			fmt.Println("Starting trim and fade...")

			duration, err := helper.DurationFromStartAndEnd(
				captionsOptions.StartTime,
				captionsOptions.EndTime,
			)
			if err != nil {
				return err
			}

			if err := whisperService.RunTrimAndFadeOnClip(
				captionsOptions.OutputDir,
				clip,
				duration,
				&captionsOptions.FadeDuration,
				captionsOptions.Verbose,
			); err != nil {
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
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.AudioPath, "audioPath", "a", "", "Path to audio")
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.VideoPath, "videoPath", "v", "", "Path to video")
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.OutputDir, "output", "o", "output", "Output directory")
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.WhisperModel, "model", "m", "base", "Transcription model (small,base,large)")
	captionCmd.PersistentFlags().BoolVar(&captionsOptions.Verbose, "verbose", false, "Verbose output")
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.StartTime, "startTime", "s", "0", "Start time")
	captionCmd.PersistentFlags().StringVarP(&captionsOptions.EndTime, "endTime", "e", "30", "End time")
	captionCmd.PersistentFlags().BoolVarP(&captionsOptions.IsDirectory, "directoryMode", "D", false, "Enabl directory mode. Both audio path and video path must be directories when using this mode")
	captionCmd.PersistentFlags().BoolVarP(&captionsOptions.NoInteract, "no-interact", "n", false, "Disable interactive mode")

	captionCmd.MarkFlagRequired("path")

	captionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(captionCmd)
}
