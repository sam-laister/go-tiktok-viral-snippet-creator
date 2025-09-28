package model

type BatchOptions struct {
	AudioPath       string
	VideoPath       string
	OutputDir       string
	WhisperModel    string
	Verbose         bool
	StartTime       string
	EndTime         string
	NoInteract      bool
	Height          int
	Width           int
	FadeDuration    int
	SkipCaptionsGen bool
	SkipVideoGen    bool
}

func NewBatchOptions(opts ...func(*BatchOptions)) *BatchOptions {
	const defaultHeight = 1920
	const defaultWidth = 1080
	const defaultFadeDuration = 5
	const defaultSkipCaptionsGen = false
	const defaultSkipVideoGen = false

	props := BatchOptions{
		Height:          defaultHeight,
		Width:           defaultWidth,
		FadeDuration:    defaultFadeDuration,
		SkipCaptionsGen: defaultSkipCaptionsGen,
		SkipVideoGen:    defaultSkipVideoGen,
	}
	for _, opt := range opts {
		opt(&props)
	}
	return &props
}
