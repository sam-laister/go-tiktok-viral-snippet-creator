package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type ScriptServiceImpl struct{}

func NewScriptServiceImpl() *ScriptServiceImpl {
	return &ScriptServiceImpl{}
}

func (w ScriptServiceImpl) Transcribe(
	inputFile,
	outputDir,
	model string,
	verbose bool,
) (*string, error) {
	app := "./scripts/generate_captions.py"

	t := time.Now().Unix()
	outputFile := fmt.Sprintf("%s/%d.ass", outputDir, t)

	args := []string{inputFile, outputFile, "--model", model}
	cmd := exec.CommandContext(context.Background(), app, args...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &outputFile, nil
}

func (w ScriptServiceImpl) BurnCaption(
	captionFile,
	videoFile,
	audioFile,
	outputDir string,
	targetWidth,
	targetHeight *int,
	verbose bool,
) (*string, error) {
	app := "./scripts/burn_captions.py"

	t := time.Now().Unix()
	outputFile := fmt.Sprintf("%s/%d-captions.mp4", outputDir, t)

	var defaultWidth = 1080
	var defaultHeight = 1920

	if targetWidth == nil {
		targetWidth = &defaultWidth
	}
	if targetHeight == nil {
		targetHeight = &defaultHeight
	}

	args := []string{
		captionFile,
		videoFile,
		audioFile,
		outputFile,
		strconv.Itoa(*targetWidth),
		strconv.Itoa(*targetHeight),
	}

	cmd := exec.CommandContext(context.Background(), app, args...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &outputFile, nil
}

func (w ScriptServiceImpl) TrimAndFade(
	inputFile,
	outputDir,
	startTime,
	duration string,
	fadeDuration *int,
	verbose bool,
) (*string, error) {
	app := "./scripts/trim_and_fade.py"

	t := time.Now().Unix()
	outputFile := fmt.Sprintf("%s/%d-trim.mp4", outputDir, t)

	var defaultFadeDuration = 3
	if fadeDuration == nil {
		fadeDuration = &defaultFadeDuration
	}

	args := []string{
		inputFile,
		outputFile,
		startTime,
		duration,
		fmt.Sprintf("--fade-duration=%d", *fadeDuration),
	}

	cmd := exec.CommandContext(context.Background(), app, args...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &outputFile, nil
}
