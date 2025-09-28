package model

type CaptionsOptions struct {
	AudioPath    string
	VideoPath    string
	Verbose      bool
	OutputDir    string
	WhisperModel string
	StartTime    string
	EndTime      string
	IsDirectory  bool
	NoInteract   bool
	Height       int
	Width        int
	FadeDuration int
}

func NewCaptionOptions(opts ...func(*CaptionsOptions)) *CaptionsOptions {
	const defaultHeight = 1920
	const defaultWidth = 1080
	const defaultFadeDuration = 5

	props := CaptionsOptions{
		Height:       defaultHeight,
		Width:        defaultWidth,
		FadeDuration: defaultFadeDuration,
	}
	for _, opt := range opts {
		opt(&props)
	}
	return &props
}
