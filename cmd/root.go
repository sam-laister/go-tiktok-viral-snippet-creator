/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var audioPath string
var videoPath string
var verbose bool
var outputDir string
var whisperModel string

var rootCmd = &cobra.Command{
	Use:   "tiktok-creator",
	Short: "A CLI tool to generate viral snippet videos",
	Long: `This project aims to create a powerful tool for combining
snippet videos with audio and auto-captioning. Inspired by the Mario
Kart Uzi videos.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&audioPath, "audioPath", "a", "", "Path to audio")
	rootCmd.PersistentFlags().StringVarP(&videoPath, "videoPath", "v", "", "Path to video")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "output", "Output directory")
	rootCmd.PersistentFlags().StringVarP(&whisperModel, "model", "m", "base", "Transcription model (small,base,large)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose output")

	rootCmd.MarkFlagRequired("path")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
