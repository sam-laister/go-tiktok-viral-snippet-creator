package helper

import (
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/service"
)

func IsValidClipQueue(clipQueue []*model.ClipDTO) bool {
	return len(clipQueue) != 0
}

func GenerateSRTCaptions(
	scriptService service.ScriptService,
	outputDir string,
	clip *model.ClipDTO,
	model string,
	verbose bool,
) error {
	srtPath, err := scriptService.Transcribe(
		clip.AudioInputPath,
		outputDir,
		model,
		verbose,
	)

	if err != nil {
		return err
	}

	clip.SRTCaptionPath = srtPath
	return nil
}

func BurnCaptions(
	scriptService service.ScriptService,
	outputDir string,
	clip *model.ClipDTO,
	targetWidth, targetHeight *int,
	verbose bool,
) error {
	finalOutput, err := scriptService.BurnCaption(
		*clip.SRTCaptionPath,
		clip.VideoInputPath,
		clip.AudioInputPath,
		outputDir,
		targetWidth,
		targetHeight,
		verbose,
	)

	if err != nil {
		return err
	}

	clip.CaptionsVideoOutputPath = finalOutput
	return nil
}

func TrimAndFade(
	scriptService service.ScriptService,
	outputDir string,
	clip *model.ClipDTO,
	startTime, duration string,
	fadeDuration *int,
	verbose bool,
) error {
	if fadeDuration == nil {
		fadeDuration = new(int)
		*fadeDuration = 5
	}

	trimmedPath, err := scriptService.TrimAndFade(
		*clip.CaptionsVideoOutputPath,
		outputDir,
		startTime,
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
