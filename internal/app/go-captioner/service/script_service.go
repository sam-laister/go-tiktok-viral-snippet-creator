package service

type ScriptService interface {
	Transcribe(inputFile, outputDir, model string, verbose bool, startTime, endTime string) (*string, error)
	BurnCaption(captionFile, videoFile, audioFile, outputDir string, targetWidth, targetHeight *int,
		startTime, endTime string, verbose bool) (*string, error)
	TrimAndFade(inputFile, outputDir, startTime, duration string, fadeDuration *int, verbose bool) (*string, error)
}
