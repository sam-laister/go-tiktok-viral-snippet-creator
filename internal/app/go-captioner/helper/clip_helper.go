package helper

import (
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/service"
)

func ClipFromInputPath(path string) *model.ClipDTO {
	return model.NewClipDTO(nil, &path, nil)
}

func InputPathArrayToClips(paths []string) []*model.ClipDTO {
	var clips []*model.ClipDTO
	for _, path := range paths {
		clips = append(clips, ClipFromInputPath(path))
	}
	return clips
}

func IsValidClipQueue(clipQueue []*model.ClipDTO) bool {
	return len(clipQueue) != 0
}

func GenerateSRTCaptions(
	whisperService service.WhisperService,
	outputDir string,
	clip *model.ClipDTO,
	verbose bool,
) error {
	srtPath, err := whisperService.Transcribe(
		*clip.InputPath,
		outputDir,
		verbose,
	)

	if err != nil {
		return err
	}

	clip.SRTCaptionPath = srtPath
	return nil
}
