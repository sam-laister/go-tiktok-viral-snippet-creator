package model

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type ClipDTO struct {
	AudioInputPath string  `json:"AudioInputPath"`
	VideoInputPath string  `json:"VideoInputPath"`
	SRTCaptionPath *string `json:"SRTCaptionPath"`
	OutputPath     *string `json:"OutputPath"`
}

func NewClipDTO(
	audioInputPath,
	videoInputPath string,
	srtCaptionPath,
	outputPath *string,
) *ClipDTO {
	return &ClipDTO{
		audioInputPath,
		videoInputPath,
		srtCaptionPath,
		outputPath,
	}
}

func (clip *ClipDTO) IsValidAudioInputPath() bool {
	return clip.AudioInputPath != ""
}

func (clip *ClipDTO) IsValidVideoInputPath() bool {
	return clip.VideoInputPath != ""
}

func (clip *ClipDTO) IsValidOutputPath() bool {
	return clip.OutputPath != nil && *clip.OutputPath != ""
}

func (clip *ClipDTO) IsValidSRTCaptionPath() bool {
	return clip.SRTCaptionPath != nil && *clip.SRTCaptionPath != ""
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
	if err := printRow("OutputPath", get(clip.OutputPath)); err != nil {
		return err
	}

	return w.Flush()
}
