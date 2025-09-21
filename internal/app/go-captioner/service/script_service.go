package service

type ScriptService interface {
	Transcribe(inputFile, outputDir, model string, verbose bool) (*string, error)
	BurnCaption(captionFile, videoFile, audioFile, outputDir string, targetWidth, targetHeight *int,
		verbose bool) (*string, error)
}
