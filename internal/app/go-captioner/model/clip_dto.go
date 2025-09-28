package model

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type ClipDTO struct {
	AudioInputPath          string  `json:"AudioInputPath"`
	VideoInputPath          string  `json:"VideoInputPath"`
	SRTCaptionPath          *string `json:"SRTCaptionPath"`
	CaptionsVideoOutputPath *string `json:"CaptionsVideoOutputPath"`
	TrimmedVideoOutputPath  *string `json:"TrimmedVideoOutputPath"`
	ID                      *int    `json:"ID"`
	Hash                    *string `json:"Hash"`
}

func NewClipDTO(
	audioInputPath,
	videoInputPath string,
	srtCaptionPath,
	captionsVideoOutputPath,
	trimmedVideoOutputPath,
	hash *string,
	id *int,

) *ClipDTO {
	return &ClipDTO{
		audioInputPath,
		videoInputPath,
		srtCaptionPath,
		captionsVideoOutputPath,
		trimmedVideoOutputPath,
		id,
		hash,
	}
}

func (clip *ClipDTO) IsValidAudioInputPath() bool {
	return clip.AudioInputPath != ""
}

func (clip *ClipDTO) IsValidVideoInputPath() bool {
	return clip.VideoInputPath != ""
}

func (clip *ClipDTO) IsValidCaptionsVideoOutputPath() bool {
	return clip.CaptionsVideoOutputPath != nil && *clip.CaptionsVideoOutputPath != ""
}

func (clip *ClipDTO) IsValidSRTCaptionPath() bool {
	return clip.SRTCaptionPath != nil && *clip.SRTCaptionPath != ""
}

func (clip *ClipDTO) IsValidTrimmedVideoOutputPath() bool {
	return clip.TrimmedVideoOutputPath != nil && *clip.TrimmedVideoOutputPath != ""
}

func (clip *ClipDTO) PrintTable() error {
	const tablePadding = 2
	w := tabwriter.NewWriter(os.Stdout, 0, 0, tablePadding, ' ', 0)

	printRow := func(field, val string) error {
		_, err := fmt.Fprintf(w, "%s\t%s\n", field, val)
		return err
	}

	get := func(p *string) string {
		if p == nil {
			return "<nil>"
		}
		return *p
	}

	if _, err := fmt.Fprintln(w, "Field\tValue"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "-----\t-----"); err != nil {
		return err
	}

	if err := printRow("SRTCaptionPath", get(clip.SRTCaptionPath)); err != nil {
		return err
	}
	if err := printRow("AudioInputPath", get(&clip.AudioInputPath)); err != nil {
		return err
	}
	if err := printRow("VideoInputPath", get(&clip.VideoInputPath)); err != nil {
		return err
	}
	if err := printRow("CaptionsVideoOutputPath", get(clip.CaptionsVideoOutputPath)); err != nil {
		return err
	}
	if err := printRow("TrimmedVideoOutputPath", get(clip.TrimmedVideoOutputPath)); err != nil {
		return err
	}

	return w.Flush()
}
