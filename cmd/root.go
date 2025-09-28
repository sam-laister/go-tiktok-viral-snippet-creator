/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

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

func init() {}
