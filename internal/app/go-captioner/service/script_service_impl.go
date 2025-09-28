package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/helper"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
)

const burnCaptionsPath = "./scripts/burn_captions.py"
const trimAndFadePath = "./scripts/trim_and_fade.py"
const generateCaptionsPath = "./scripts/generate_captions.py"

type ScriptServiceImpl struct{}

func NewScriptServiceImpl() *ScriptServiceImpl {
	return &ScriptServiceImpl{}
}

func (w ScriptServiceImpl) RunGenerateSRTCaptionsOnClip(
	outputDir string,
	clip *model.ClipDTO,
	model string,
	startTime,
	endTime string,
	verbose bool,
) error {
	srtPath, err := w.Transcribe(
		clip.AudioInputPath,
		outputDir,
		model,
		verbose,
		startTime,
		endTime,
	)

	if err != nil {
		return err
	}

	clip.SRTCaptionPath = srtPath
	return nil
}

func (w ScriptServiceImpl) RunBurnCaptionsOnClip(
	outputDir string,
	clip *model.ClipDTO,
	targetWidth, targetHeight *int,
	startTime, endTime string,
	verbose bool,
) error {
	if clip.SRTCaptionPath == nil {
		return errors.New("no captions path provided")
	}

	finalOutput, err := w.BurnCaption(
		*clip.SRTCaptionPath,
		clip.VideoInputPath,
		clip.AudioInputPath,
		outputDir,
		targetWidth,
		targetHeight,
		startTime,
		endTime,
		verbose,
	)

	if err != nil {
		return err
	}

	clip.CaptionsVideoOutputPath = finalOutput
	return nil
}

func (w ScriptServiceImpl) RunTrimAndFadeOnClip(
	outputDir string,
	clip *model.ClipDTO,
	duration string,
	fadeDuration *int,
	verbose bool,
) error {
	if fadeDuration == nil {
		fadeDuration = new(int)
		*fadeDuration = 5
	}

	trimmedPath, err := w.TrimAndFade(
		*clip.CaptionsVideoOutputPath,
		outputDir,
		duration,
		fadeDuration,
		verbose,
	)

	if err != nil {
		return err
	}

	clip.TrimmedVideoOutputPath = trimmedPath
	return nil
}

func (w ScriptServiceImpl) Transcribe(
	inputFile,
	outputDir,
	model string,
	verbose bool,
	startTime,
	endTime string,
) (*string, error) {
	t := time.Now().Unix()
	outputFile := fmt.Sprintf("%s/%d.ass", outputDir, t)

	args := []string{inputFile, outputFile, "--model", model, "--start", startTime, "--end", endTime}
	cmd := exec.CommandContext(context.Background(), generateCaptionsPath, args...)

	fmt.Println("Running: ", helper.GetCommandPrintable(cmd))

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
	startTime,
	endTime string,
	verbose bool,
) (*string, error) {
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
		"--start",
		startTime,
		"--end",
		endTime,
	}

	cmd := exec.CommandContext(context.Background(), burnCaptionsPath, args...)

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
	duration string,
	fadeDuration *int,
	verbose bool,
) (*string, error) {
	t := time.Now().Unix()
	outputFile := fmt.Sprintf("%s/%d-trim.mp4", outputDir, t)

	var defaultFadeDuration = 3
	if fadeDuration == nil {
		fadeDuration = &defaultFadeDuration
	}

	args := []string{
		inputFile,
		outputFile,
		"0", // Vid is pre cropped at this point
		duration,
		fmt.Sprintf("--fade-duration=%d", *fadeDuration),
	}

	cmd := exec.CommandContext(context.Background(), trimAndFadePath, args...)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &outputFile, nil
}
