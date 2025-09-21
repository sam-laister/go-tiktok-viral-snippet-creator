package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type WhisperServiceImpl struct{}

func NewWhisperServiceImpl() *WhisperServiceImpl {
	return &WhisperServiceImpl{}
}

func (w WhisperServiceImpl) Transcribe(inputFile, outputDir string, verbose bool) (*string, error) {
	app := "./scripts/generate_captions.py"

	t := time.Now()
	outputFile := fmt.Sprintf("%s/%s.srt", outputDir, t.Format(time.DateTime))

	arg0 := inputFile
	arg1 := outputFile

	cmd := exec.Command(app, arg0, arg1)

	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		log.Println("ERROR Running:", app, arg0, arg1, arg1)
		return nil, err
	}

	return &outputFile, nil
}
