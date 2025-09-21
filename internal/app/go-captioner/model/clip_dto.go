package model

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type ClipDTO struct {
	SRTCaptionPath *string `json:"SRTCaptionPath"`
	InputPath      *string `json:"InputPath"`
	OutputPath     *string `json:"OutputPath"`
}

func NewClipDTO(
	srtCaptionPath,
	inputPath,
	outputPath *string,
) *ClipDTO {
	return &ClipDTO{
		srtCaptionPath,
		inputPath,
		outputPath,
	}
}

func (clip *ClipDTO) IsValidInputPath() bool {
	return clip.InputPath != nil || *clip.InputPath != ""
}

func (clip *ClipDTO) IsValidOutputPath() bool {
	return clip.OutputPath != nil && *clip.OutputPath != ""
}

func (clip *ClipDTO) IsValidSRTCaptionPath() bool {
	return clip.SRTCaptionPath != nil && *clip.SRTCaptionPath != ""
}

func (clip *ClipDTO) PrintTable() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

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
	if err := printRow("InputPath", get(clip.InputPath)); err != nil {
		return err
	}
	if err := printRow("OutputPath", get(clip.OutputPath)); err != nil {
		return err
	}

	return w.Flush()
}
