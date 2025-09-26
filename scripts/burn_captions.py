#!/usr/bin/env python3
"""
Burn ASS captions into a video using ffmpeg.
Usage:
    python burn_captions.py captions_file.ass video_file audio_file output_file [target_width] [target_height]
        [--start 0] [--end 30]
"""

import argparse
import os
import sys
import random
import ffmpeg

def main():
    parser = argparse.ArgumentParser(description="Burn ASS captions into a video using ffmpeg")
    parser.add_argument("captions_file", help="Path to the captions .ass file")
    parser.add_argument("video_file", help="Path to the video file")
    parser.add_argument("audio_file", help="Path to the audio file")
    parser.add_argument("output_file", help="Path to write output video to")
    parser.add_argument("target_width", nargs="?", type=int, default=1080)
    parser.add_argument("target_height", nargs="?", type=int, default=1920)
    parser.add_argument("--start", type=float, default=0.0, help="Start time in seconds (default: 0)")
    parser.add_argument("--end", type=float, default=30.0, help="End time in seconds (default: 30)")
    args = parser.parse_args()

    if not os.path.isfile(args.captions_file):
        sys.exit(f"Error: {args.captions_file} does not exist.")
    if not os.path.isfile(args.video_file):
        sys.exit(f"Error: {args.video_file} does not exist.")
    if not os.path.isfile(args.audio_file):
        sys.exit(f"Error: {args.audio_file} does not exist.")

    # Validate times and compute duration
    if args.end <= args.start:
        sys.exit("Error: --end must be greater than --start")
    clip_duration = args.end - args.start

    # Probe duration and choose random start within bounds
    meta = ffmpeg.probe(args.video_file)
    total_duration = float(meta.get("format", {}).get("duration", 0.0))
    if not total_duration:
        vstream = next(s for s in meta["streams"] if s["codec_type"] == "video")
        total_duration = float(vstream.get("duration", 0.0))
    if total_duration and clip_duration > total_duration:
        clip_duration = total_duration
    max_start = max(0.0, total_duration - clip_duration)
    trim_start = round(random.uniform(0.0, max_start), 3)

    print(f"Clip duration: {clip_duration}s")
    print(f"Random start: {trim_start}s of total {total_duration}s")

    target_aspect = args.target_width / args.target_height

    # Crop video to target aspect ratio
    probe = ffmpeg.probe(args.video_file)
    vinfo = next(s for s in probe["streams"] if s["codec_type"] == "video")
    in_w, in_h = int(vinfo["width"]), int(vinfo["height"])
    input_aspect = in_w / in_h
    if input_aspect > target_aspect:          # too wide
        crop_h = in_h
        crop_w = int(crop_h * target_aspect)
        crop_x, crop_y = (in_w - crop_w)//2, 0
    else:                                     # too tall
        crop_w = in_w
        crop_h = int(crop_w / target_aspect)
        crop_x, crop_y = 0, (in_h - crop_h)//2

    video = (
        ffmpeg
        .input(args.video_file, ss=trim_start, t=clip_duration)
        .video
        .filter("crop", crop_w, crop_h, crop_x, crop_y)
        .filter("scale", args.target_width, args.target_height)
        .filter("ass", args.captions_file)
    )
    audio = (
        ffmpeg
        .input(args.audio_file, ss=args.start, t=clip_duration)
        .audio
    )

    (
        ffmpeg
        .output(video, audio, args.output_file, acodec='aac',
                audio_bitrate='320k', shortest=None)
        .overwrite_output()
        .run()
    )
    print("Done!")

if __name__ == "__main__":
    main()
