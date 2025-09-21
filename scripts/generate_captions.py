#!/usr/bin/env python3
"""
Generate captions (SRT) locally using the openai-whisper library.

Usage:
    python generate_captions_local.py input_file output_dir [--model base]
"""

import argparse
import datetime
import os
import sys
import whisper

def main():
    parser = argparse.ArgumentParser(description="Create captions locally using Whisper")
    parser.add_argument("input_file", help="Path to the audio or video file")
    parser.add_argument("output_file", help="Path to write output srt to")
    parser.add_argument("--model", default="base",
                        help="Whisper model size: tiny, base, small, medium, large")
    args = parser.parse_args()

    input_path = args.input_file
    output_file = args.output_file
    model_size = args.model

    if not os.path.isfile(input_path):
        sys.exit(f"Error: {input_path} does not exist.")
    os.makedirs(output_file, exist_ok=True)

    print(f"Loading Whisper model '{model_size}' … (first time may take a while)")
    model = whisper.load_model(model_size)

    print(f"Transcribing {input_path} …")
    # transcribe() automatically converts video to audio if needed
    result = model.transcribe(input_path, verbose=True)

    print(f"Writing captions to {output_file}")

    writer = whisper.utils.get_writer("srt", output_file)
    writer(result, output_file)

    print("Done!")

if __name__ == "__main__":
    main()
