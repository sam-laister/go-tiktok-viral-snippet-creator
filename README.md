## Tiktok snippet Generator

```
This project aims to create a powerful tool for combining
snippet videos with audio and auto-captioning. Inspired by the Mario
Kart Uzi videos.

Usage:
  tiktok-creator [command]

Available Commands:
  caption     Generates captions for an audio file.
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -a, --audioPath string   Path to audio
  -h, --help               help for tiktok-creator
  -m, --model string       Transcription model (small,base,large) (default "base")
  -o, --output string      Output directory (default "output")
  -t, --toggle             Help message for toggle
      --verbose            Verbose output
  -v, --videoPath string   Path to video

Use "tiktok-creator [command] --help" for more information about a command.
```

### Dependencies

Ensure Go is installed. https://go.dev/learn/

`python -m pip install -r requirements.txt`

### Project Example

`go run main.go caption -a sample/audio.mp3 -o output -v sample/bg.mp4 -m large --verbose`
