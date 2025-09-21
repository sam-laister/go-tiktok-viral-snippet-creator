package service

type WhisperService interface {
	Transcribe(inputFile, outputDir string, verbose bool) (*string, error)
}
