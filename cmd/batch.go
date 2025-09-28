/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/helper"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/repository"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/service"
	"github.com/spf13/cobra"
)

var batchOptions = model.NewBatchOptions()

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Batch generate captions using an SQLite database to track progress.",
	Long:  `Batch generate captions using an SQLite database to track progress.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := helper.GetDB()
		if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
		}
		defer client.Close()

		clipRepository := repository.NewClipRepository(client)

		whisperService := service.NewScriptServiceImpl()
		clipService := service.NewClipServiceImpl(clipRepository)

		fmt.Println("Verbose: ", batchOptions.Verbose)

		if !helper.IsDirectory(batchOptions.AudioPath) || !helper.IsDirectory(batchOptions.VideoPath) {
			return errors.New(fmt.Sprintf(
				"either %s or %s is not a directory",
				batchOptions.AudioPath,
				batchOptions.VideoPath,
			))
		}

		audios, err := helper.GetFilesInDirectory(batchOptions.AudioPath)
		if err != nil {
			return err
		}

		videos, err := helper.GetFilesInDirectory(batchOptions.VideoPath)
		if err != nil {
			return err
		}

		for _, audioPath := range audios {
			audioHash, err := helper.GetFilehash(audioPath)
			if err != nil {
				fmt.Println(fmt.Sprintf("Failed calculating hash for file %s %s", audioPath, err.Error()))
				continue
			}

			// Database entry
			clipDTO, err := clipService.GetOrCreateWithHash(
				context.Background(),
				audioHash,
				audioPath,
				videos[rand.Intn(len(videos))],
			)
			if err != nil {
				fmt.Println(fmt.Sprintf("Failed creating clip for file %s %s", audioPath, err.Error()))
				continue
			}

			// Captions Gen
			if !batchOptions.SkipCaptionsGen &&
				(clipDTO.SRTCaptionPath == nil ||
					!helper.Exists(*clipDTO.SRTCaptionPath)) {
				if err := whisperService.RunGenerateSRTCaptionsOnClip(
					batchOptions.OutputDir,
					clipDTO,
					batchOptions.WhisperModel,
					batchOptions.StartTime,
					batchOptions.EndTime,
					batchOptions.Verbose,
				); err != nil {
					fmt.Println(fmt.Sprintf("Failed generating captions for file %s %s", audioPath, err.Error()))
					continue
				}
				if err = clipService.Update(context.Background(), clipDTO); err != nil {
					fmt.Println(fmt.Sprintf("Failed updating clip for file %s %s", audioPath, err.Error()))
					continue
				}
			} else {
				fmt.Println(fmt.Sprintf("Skipping caption generation for file %s...", audioPath))
			}

			if err := clipDTO.PrintTable(); err != nil {
				return err
			}

			if !batchOptions.NoInteract {
				if _, err := helper.WaitForOptionalEdits(); err != nil {
					return err
				}
			}

			// Raw Video Gen
			if !batchOptions.SkipVideoGen &&
				(clipDTO.CaptionsVideoOutputPath == nil ||
					!helper.Exists(*clipDTO.CaptionsVideoOutputPath)) {
				fmt.Println("Starting burn...")
				if err := whisperService.RunBurnCaptionsOnClip(
					batchOptions.OutputDir,
					clipDTO,
					&batchOptions.Width,
					&batchOptions.Height,
					batchOptions.StartTime,
					batchOptions.EndTime,
					batchOptions.Verbose,
				); err != nil {
					fmt.Println(fmt.Sprintf("Failed burning captions for file %s %s", audioPath, err.Error()))
					continue
				}

				if err = clipService.Update(context.Background(), clipDTO); err != nil {
					fmt.Println(fmt.Sprintf("Failed generating video clip for file %s %s", audioPath, err.Error()))
					continue
				}
			} else {
				fmt.Println("Burn already exists, skipping...")
			}

			// Final Video gen
			if !batchOptions.SkipVideoGen &&
				(clipDTO.TrimmedVideoOutputPath == nil ||
					!helper.Exists(*clipDTO.TrimmedVideoOutputPath)) {
				fmt.Println("Starting trim and fade...")
				duration, err := helper.DurationFromStartAndEnd(
					batchOptions.StartTime,
					batchOptions.EndTime,
				)
				if err != nil {
					return err
				}
				if err := whisperService.RunTrimAndFadeOnClip(
					batchOptions.OutputDir,
					clipDTO,
					duration,
					&batchOptions.FadeDuration,
					batchOptions.Verbose,
				); err != nil {
					return err
				}

				if err = clipService.Update(context.Background(), clipDTO); err != nil {
					fmt.Println(fmt.Sprintf("Failed trimming video clip for file %s %s", audioPath, err.Error()))
					continue
				}
			} else {
				fmt.Println("Trimmed output already exists, skipping...")
			}

			if err := clipDTO.PrintTable(); err != nil {
				return err
			}

		}

		return nil
	},
}

func init() {
	batchCmd.PersistentFlags().StringVarP(&batchOptions.AudioPath, "audioPath", "a", "", "Path to audio")
	batchCmd.PersistentFlags().StringVarP(&batchOptions.VideoPath, "videoPath", "v", "", "Path to video")
	batchCmd.PersistentFlags().StringVarP(&batchOptions.OutputDir, "output", "o", "output", "Output directory")
	batchCmd.PersistentFlags().StringVarP(&batchOptions.WhisperModel, "model", "m", "base", "Transcription model (small,base,large)")
	batchCmd.PersistentFlags().BoolVar(&batchOptions.Verbose, "verbose", false, "Verbose output")
	batchCmd.PersistentFlags().StringVarP(&batchOptions.StartTime, "startTime", "s", "0", "Start time")
	batchCmd.PersistentFlags().StringVarP(&batchOptions.EndTime, "endTime", "e", "30", "End time")
	batchCmd.PersistentFlags().BoolVarP(&batchOptions.NoInteract, "no-interact", "n", false, "Disable interactive mode")
	batchCmd.PersistentFlags().BoolVar(&batchOptions.SkipVideoGen, "skip-video-gen", false, "Skip Video Generation")
	batchCmd.PersistentFlags().BoolVar(&batchOptions.SkipCaptionsGen, "skip-captions-gen", false, "Skip Captions Generation")

	batchCmd.MarkFlagRequired("audioPath")
	batchCmd.MarkFlagRequired("videoPath")
	batchCmd.MarkFlagRequired("output")

	batchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(batchCmd)
}
