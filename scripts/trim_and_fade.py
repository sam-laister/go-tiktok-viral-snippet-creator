#!/usr/bin/env python3
"""
Trim video and add fade-out effects to both audio and video.
Usage:
    python trim_and_fade.py input_file output_file start_time duration [--fade-duration 3]
"""

import argparse
import os
import sys
import ffmpeg

def trim_and_fade(input_file, output_file, start_time, duration, fade_duration=3):
    """Trim video and add fade-out effects."""

    # Calculate end time
    end_time = start_time + duration

    print(f"Trimming video to {duration} from {start_time} to {end_time}")
    print(f"Adding {fade_duration}s fade-out effect")

    # Get video info
    probe = ffmpeg.probe(input_file)
    video_stream = next(s for s in probe["streams"] if s["codec_type"] == "video")
    audio_stream = next((s for s in probe["streams"] if s["codec_type"] == "audio"), None)

    input_width = int(video_stream["width"])
    input_height = int(video_stream["height"])

    # Calculate fade start time (fade_duration seconds before end)
    fade_start_time = max(0, duration - fade_duration)

    print(f"Video dimensions: {input_width}x{input_height}")
    print(f"Fade starts at {fade_start_time} of trimmed video")

    # Process video with fade-out
    video = (
        ffmpeg.input(input_file, ss=start_time, t=duration)
        .filter('scale', input_width, input_height)  # Ensure consistent scaling
        .filter('fade', t='out', st=fade_start_time, d=fade_duration)
    )

    # Process audio if present
    if audio_stream:
        audio = (
            ffmpeg.input(input_file, ss=start_time, t=duration)
            .filter('afade', t='out', st=fade_start_time, d=fade_duration)
        )

        # Combine video and audio
        out = ffmpeg.output(
            video, audio, output_file,
            vcodec='libx264',
            acodec='aac',
            audio_bitrate='320k',
            preset='fast'
        )
    else:
        # Video only
        out = ffmpeg.output(
            video, output_file,
            vcodec='libx264',
            preset='fast'
        )

    # Run ffmpeg
    out.overwrite_output().run(quiet=False)

def main():
    parser = argparse.ArgumentParser(description="Trim video and add fade-out effects")
    parser.add_argument("input_file", help="Path to input video file")
    parser.add_argument("output_file", help="Path to output video file")
    parser.add_argument("start_time", type=int, help="Start time in seconds (Default: 0)", default=0)
    parser.add_argument("duration", type=int, help="Duration in seconds (Default: 60)", default=60)
    parser.add_argument("--fade-duration", type=int, default=3,
                       help="Fade-out duration in seconds (default: 3)")
    args = parser.parse_args()

    # Validate input file
    if not os.path.isfile(args.input_file):
        sys.exit(f"Error: {args.input_file} does not exist.")

    # Create output directory if needed
    os.makedirs(os.path.dirname(args.output_file) or ".", exist_ok=True)

    try:
        print(f"Input file: {args.input_file}")
        print(f"Output file: {args.output_file}")
        print(f"Start time: {args.start_time}s")
        print(f"Duration: {args.duration}s")
        print(f"Fade duration: {args.fade_duration}s")

        # Validate fade duration
        if args.fade_duration >= args.duration:
            sys.exit(f"Error: Fade duration ({args.fade_duration}s) must be less than video duration ({args.duration}s)")

        # Process the video
        trim_and_fade(args.input_file, args.output_file, args.start_time, args.duration, args.fade_duration)

        print("Done!")

    except ValueError as e:
        sys.exit(f"Error: {e}")
    except Exception as e:
        sys.exit(f"Error processing video: {e}")

if __name__ == "__main__":
    main()
