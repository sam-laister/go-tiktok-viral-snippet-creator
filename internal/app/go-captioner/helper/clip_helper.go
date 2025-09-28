package helper

import (
	"github.com/sam-laister/tiktok-creator/ent"
	"github.com/sam-laister/tiktok-creator/internal/app/go-captioner/model"
)

func IsValidClipQueue(clipQueue []*model.ClipDTO) bool {
	return len(clipQueue) != 0
}

func ClipToDTO(c *ent.Clip) *model.ClipDTO {
	var id *int
	if c.ID != 0 {
		id = &c.ID
	}

	var hash *string
	if c.Hash != "" {
		hash = &c.Hash
	}

	return &model.ClipDTO{
		AudioInputPath:          c.AudioPath,
		VideoInputPath:          c.VideoPath,
		SRTCaptionPath:          c.GenCaptionsPath,
		CaptionsVideoOutputPath: c.GenRawVideoPath,
		TrimmedVideoOutputPath:  c.GenTrimmedVideoPath,
		ID:                      id,
		Hash:                    hash,
	}
}

func DTOToClip(dto *model.ClipDTO) *ent.Clip {
	var id int
	if dto.ID != nil {
		id = *dto.ID
	}

	var hash string
	if dto.Hash != nil {
		hash = *dto.Hash
	}

	return &ent.Clip{
		ID:                  id,
		Hash:                hash,
		AudioPath:           dto.AudioInputPath,
		VideoPath:           dto.VideoInputPath,
		GenCaptionsPath:     dto.SRTCaptionPath,
		GenRawVideoPath:     dto.CaptionsVideoOutputPath,
		GenTrimmedVideoPath: dto.TrimmedVideoOutputPath,
	}
}
