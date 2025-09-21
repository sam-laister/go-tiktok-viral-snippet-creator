/*
Copyright © 2025 Sam Laister <laister.sam@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var path string
var isDirectory bool
var verbose bool
var outputDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiktok-creator",
	Short: "A CLI tool to generate viral snippet videos",
	Long: `This project aims to create a powerful tool for combining 
snippet videos with audio and auto-captioning. Inspired by the Mario 
Kart Uzi videos.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Path to audio")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "output", "Output directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.MarkFlagRequired("path")

	rootCmd.PersistentFlags().BoolVarP(&isDirectory, "directory", "d", false, "Specifies if path is a directory")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
