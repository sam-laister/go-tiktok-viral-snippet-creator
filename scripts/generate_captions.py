#!/usr/bin/env python3
"""
Generate ASS subtitles that fill a 1080x1920 screen in a grid.
Each word occupies a cell; when the grid is full the screen clears and
the next page begins.

Usage:
    python generate_grid_ass.py input_file output_file.ass [--model base]
"""

import argparse
import os
import sys
import whisper

# ==== Layout settings for 1080x1920 portrait video ====
VIDEO_W, VIDEO_H = 1080, 1920
COLS, ROWS = 2, 4
MARGIN_X, MARGIN_Y = 60, 240         # top/left padding in pixels
CELL_W = (VIDEO_W - 2 * MARGIN_X) // COLS
CELL_H = (VIDEO_H - 2 * MARGIN_Y) // ROWS
FONT_SIZE = 200                      # adjust to taste

def write_grid_ass(result, out_path: str):
    header = (
        f"[Script Info]\nPlayResX: {VIDEO_W}\nPlayResY: {VIDEO_H}\nScriptType: v4.00+\n\n"
        "[V4+ Styles]\n"
        "Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, "
        "OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, "
        "ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, "
        "Alignment, MarginL, MarginR, MarginV, Encoding\n"
        f"Style: Default,Ubuntu,{FONT_SIZE},&H00000000,&H000000FF,&H00FFFFFF,&H64000000,"
        "0,0,0,0,100,100,0,0,1,2,0,5,10,10,10,1\n\n"
        "[Events]\n"
        "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"
    )

    words = []
    for seg in result["segments"]:
        words.extend(seg.get("words", []))

    page_size = COLS * ROWS
    events = []

    def to_ass_time(t):
        h = int(t // 3600)
        m = int((t % 3600) // 60)
        s = t % 60
        return f"{h:d}:{m:02d}:{s:05.2f}"

    for page_idx in range(0, len(words), page_size):
        page_words = words[page_idx: page_idx + page_size]
        # End time for the entire page = start of next page's first word or last word's end
        if page_idx + page_size < len(words):
            page_end = words[page_idx + page_size]["start"]
        else:
            page_end = page_words[-1]["end"]

        for i, w in enumerate(page_words):
            col = i % COLS
            row = i // COLS
            x = MARGIN_X + col * CELL_W + CELL_W // 2
            y = MARGIN_Y + row * CELL_H + CELL_H // 2
            text = w["word"].strip()
            if not text:
                continue
            events.append(
                f"Dialogue: 0,{to_ass_time(w['start'])},{to_ass_time(page_end)},"
                f"Default,,0,0,0,,{{\\pos({x},{y})}}{text}"
            )

    with open(out_path, "w", encoding="utf-8") as f:
        f.write(header + "\n".join(events))



def main():
    parser = argparse.ArgumentParser(description="Generate grid-layout ASS subtitles for 1080x1920 video")
    parser.add_argument("input_file", help="Path to audio/video file")
    parser.add_argument("output_file", help="Output .ass file path")
    parser.add_argument("--model", default="base", help="Whisper model size")
    args = parser.parse_args()

    if not os.path.isfile(args.input_file):
        sys.exit(f"Error: {args.input_file} does not exist.")
    os.makedirs(os.path.dirname(args.output_file) or ".", exist_ok=True)

    print(f"Loading Whisper model '{args.model}'…")
    model = whisper.load_model(args.model)

    print("Transcribing with word timestamps…")
    result = model.transcribe(args.input_file, word_timestamps=True, verbose=True)

    print(f"Writing ASS grid subtitles to {args.output_file}")
    write_grid_ass(result, args.output_file)
    print("Done!")

if __name__ == "__main__":
    main()
