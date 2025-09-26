#!/usr/bin/env python3
"""
Generate ASS subtitles with dynamic text flow based on character limits.
Text flows naturally with word wrapping, creating a more readable layout.

Usage:
    python generate_captions.py input_file output_file.ass [--model base]
        [--max-chars 50] [--font-size 200] [--margins 60,240]
        [--start 0] [--end 30]
"""

import argparse
import os
import sys
import tempfile
import ffmpeg
import whisper

# ==== Layout settings for 1080x1920 portrait video ====
VIDEO_W, VIDEO_H = 1080, 1920
DEFAULT_MAX_CHARS = 12               # characters per line
DEFAULT_FONT_SIZE = 200
DEFAULT_MARGIN_X, DEFAULT_MARGIN_Y = 60, 240

def wrap_text_to_lines(words, max_chars_per_line):
    """Wrap words into lines based on character count."""
    lines = []
    current_line = []
    current_chars = 0

    for word_info in words:
        word = word_info["word"].strip()
        if not word:
            continue

        # Check if adding this word would exceed the limit
        if current_chars + len(word) + (1 if current_line else 0) > max_chars_per_line:
            if current_line:  # Only create new line if we have content
                lines.append(current_line)
                current_line = []
                current_chars = 0

        current_line.append(word_info)
        current_chars += len(word) + (1 if current_line else 0)

    # Add the last line if it has content
    if current_line:
        lines.append(current_line)

    return lines

def write_dynamic_ass(result, out_path: str, max_chars_per_line: int, font_size: int, margin_x: int, margin_y: int):
    """Generate ASS subtitles with dynamic text flow and individual word timing."""
    header = (
        f"[Script Info]\nPlayResX: {VIDEO_W}\nPlayResY: {VIDEO_H}\nScriptType: v4.00+\n\n"
        "[V4+ Styles]\n"
        "Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, "
        "OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, "
        "ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, "
        "Alignment, MarginL, MarginR, MarginV, Encoding\n"
        f"Style: Default,Ubuntu,{font_size},&H00000000,&H000000FF,&H00FFFFFF,&H64000000,"
        "0,0,0,0,100,100,0,0,1,3,2,5,10,10,10,1\n\n"
        "[Events]\n"
        "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"
    )

    words = []
    for seg in result["segments"]:
        words.extend(seg.get("words", []))

    # Wrap text into lines
    lines = wrap_text_to_lines(words, max_chars_per_line)

    events = []

    def to_ass_time(t):
        h = int(t // 3600)
        m = int((t % 3600) // 60)
        s = t % 60
        return f"{h:d}:{m:02d}:{s:05.2f}"

    # Calculate layout parameters
    line_height = int(font_size * 1.4)  # 1.4x font size for line spacing
    max_lines_per_page = (VIDEO_H - 2 * margin_y) // line_height

    print(f"Max lines per page: {max_lines_per_page}")
    print(f"Lines: {lines}")

    # Process lines in pages
    for page_idx in range(0, len(lines), max_lines_per_page):
        page_lines = lines[page_idx:page_idx + max_lines_per_page]
        print(f"Page lines: {page_lines}")
        # Calculate timing for the entire page
        if page_lines:
            page_start = page_lines[0][0]["start"]
            if page_idx + max_lines_per_page < len(lines):
                page_end = lines[page_idx + max_lines_per_page][0]["start"]
            else:
                page_end = page_lines[-1][-1]["end"]
            print(f"Page start: {page_start}")
            print(f"Page end: {page_end}")
            # Create individual word events that persist until page end
            for line_idx, line_words in enumerate(page_lines):
                y_pos = margin_y + line_idx * line_height + line_height // 2

                # Calculate line width and starting position for centering
                total_line_text = " ".join(w["word"].strip() for w in line_words if w["word"].strip())
                line_width_chars = len(total_line_text)

                # More conservative character width estimate for monospace-like behavior
                char_width = font_size * 0.5  # More conservative estimate
                line_width_pixels = line_width_chars * char_width

                # Ensure line doesn't exceed screen width with margins
                available_width = VIDEO_W - 2 * margin_x
                if line_width_pixels > available_width:
                    # Scale down character width to fit
                    char_width = available_width / line_width_chars
                    line_width_pixels = available_width

                line_start_x = (VIDEO_W - line_width_pixels) // 2
                line_start_x = max(margin_x, line_start_x)  # Ensure we don't go negative

                print(f"Line: '{total_line_text}' - width: {line_width_pixels:.0f}px, start_x: {line_start_x:.0f}")

                # Create individual word events
                current_x = line_start_x
                for word_idx, word_info in enumerate(line_words):
                    word = word_info["word"].strip()
                    if not word:
                        continue

                    # Add space before word (except for first word in line)
                    if word_idx > 0:
                        current_x += char_width  # Add space width

                    word_width = len(word) * char_width
                    word_center_x = current_x + word_width // 2

                    print(f"Word: '{word}' at position {word_center_x:.0f},{y_pos}")

                    # Create event for this individual word appearing and persisting until page end
                    events.append(
                        f"Dialogue: 0,{to_ass_time(word_info['start'])},{to_ass_time(page_end)},"
                        f"Default,,0,0,0,,{{\\pos({word_center_x:.0f},{y_pos})}}{word}"
                    )

                    # Move to next word position
                    current_x += word_width
    with open(out_path, "w", encoding="utf-8") as f:
        f.write(header + "\n".join(events))



def main():
    parser = argparse.ArgumentParser(description="Generate dynamic ASS subtitles for 1080x1920 video")
    parser.add_argument("input_file", help="Path to audio/video file")
    parser.add_argument("output_file", help="Output .ass file path")
    parser.add_argument("--model", default="base", help="Whisper model size")
    parser.add_argument("--max-chars", type=int, default=DEFAULT_MAX_CHARS,
                       help=f"Maximum characters per line (default: {DEFAULT_MAX_CHARS})")
    parser.add_argument("--font-size", type=int, default=DEFAULT_FONT_SIZE,
                       help=f"Font size in pixels (default: {DEFAULT_FONT_SIZE})")
    parser.add_argument("--margins", default=f"{DEFAULT_MARGIN_X},{DEFAULT_MARGIN_Y}",
                       help=f"Margins as 'x,y' in pixels (default: {DEFAULT_MARGIN_X},{DEFAULT_MARGIN_Y})")
    parser.add_argument("--start", type=float, default=0.0,
                       help="Start time in seconds to begin transcription (default: 0)")
    parser.add_argument("--end", type=float, default=30.0,
                       help="End time in seconds for transcription (default: 30)")
    args = parser.parse_args()

    if not os.path.isfile(args.input_file):
        sys.exit(f"Error: {args.input_file} does not exist.")
    os.makedirs(os.path.dirname(args.output_file) or ".", exist_ok=True)

    # Parse margins
    try:
        margin_x, margin_y = map(int, args.margins.split(','))
    except ValueError:
        sys.exit("Error: Margins must be in format 'x,y' (e.g., '60,240')")

    # Determine segment to process
    if args.end <= args.start:
        sys.exit("Error: --end must be greater than --start")
    segment_duration = args.end - args.start

    # Create a temporary trimmed audio file for faster transcription and correctness
    tmp_dir = tempfile.mkdtemp(prefix="capseg_")
    tmp_audio = os.path.join(tmp_dir, "segment.wav")

    print(f"Trimming input to segment [{args.start:.2f}s, {args.end:.2f}s) → duration {segment_duration:.2f}s")
    (
        ffmpeg
        .input(args.input_file, ss=args.start, t=segment_duration)
        .output(tmp_audio, acodec='pcm_s16le', ac=1, ar='16000')
        .overwrite_output()
        .run(quiet=False)
    )

    print(f"Loading Whisper model '{args.model}'…")
    model = whisper.load_model(args.model)

    print("Transcribing trimmed segment with word timestamps…")
    result = model.transcribe(tmp_audio, word_timestamps=True, verbose=True)

    print(f"Writing dynamic ASS subtitles to {args.output_file}")
    print(f"Settings: max_chars={args.max_chars}, font_size={args.font_size}, margins={margin_x},{margin_y}")
    write_dynamic_ass(result, args.output_file, args.max_chars, args.font_size, margin_x, margin_y)

    # Cleanup
    try:
        os.remove(tmp_audio)
        os.rmdir(tmp_dir)
    except Exception:
        pass
    print("Done!")

if __name__ == "__main__":
    main()
