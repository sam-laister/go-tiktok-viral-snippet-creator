#!/usr/bin/env python3
"""
Burn captions (SRT) into a video using ffmpeg

Usage:
    python burn_captions.py captions_file video_file output_dir
"""

import argparse
import os
import sys
import ffmpeg

def main():
    parser = argparse.ArgumentParser(description="Burn captions (SRT) into a video using ffmpeg")
    parser.add_argument("captions_file", help="Path to the captions srt file")
    parser.add_argument("video_file", help="Path to the video file to")
    parser.add_argument("audio_file", help="Path to the audio file to")
    parser.add_argument("output_file", help="Path to write output video to")
    parser.add_argument("target_width", help="Target width for video", nargs='?', default=1080, type=int)
    parser.add_argument("target_height", help="Target height for video", nargs='?', default=1920, type=int)
    args = parser.parse_args()

    captions_path = args.captions_file
    video_path = args.video_file
    audio_path = args.audio_file
    output_path = args.output_file
    target_width = args.target_width
    target_height = args.target_height

    if not os.path.isfile(captions_path):
        sys.exit(f"Error: {captions_path} does not exist.")

    if not os.path.isfile(video_path):
        sys.exit(f"Error: {video_path} does not exist.")

    if not os.path.isfile(audio_path):
        sys.exit(f"Error: {audio_path} does not exist.")

    subtitle_style = (
        "FontName=Arial,"      # font
        "FontSize=42,"         # huge font size
        "PrimaryColour=&HFFFFFF,"  # white text
        "OutlineColour=&H000000,"  # black outline
        "Outline=2,"           # thickness of outline
        "BorderStyle=1,"       # 1=outline+shadow
        "Alignment=2,"         # center bottom (1=bottom-left, 2=bottom-center, 3=bottom-right, 5=middle-center)
        "MarginV=20"           # vertical margin from bottom
    )

    target_aspect = target_width / target_height

    video = ffmpeg.input(video_path)
    audio = ffmpeg.input(audio_path)

    # Probe input video
    probe = ffmpeg.probe(video_path)
    video_info = next(stream for stream in probe['streams'] if stream['codec_type'] == 'video')
    in_width = int(video_info['width'])
    in_height = int(video_info['height'])
    input_aspect = in_width / in_height

    # Calculate crop to 9:16 aspect ratio
    if input_aspect > target_aspect:
        # Input wider than 9:16 -> crop width
        crop_height = in_height
        crop_width = int(crop_height * target_aspect)
        crop_x = (in_width - crop_width) // 2
        crop_y = 0
    else:
        # Input taller than 9:16 -> crop height
        crop_width = in_width
        crop_height = int(crop_width / target_aspect)
        crop_x = 0
        crop_y = (in_height - crop_height) // 2

    command = (
        ffmpeg
        .output(
            video
            .filter('crop', crop_width, crop_height, crop_x, crop_y)
            .filter('scale', target_width, target_height)
            .filter('subtitles', captions_path, force_style=subtitle_style),
            audio,
            output_path,
            vcodec='libx264', acodec='aac'
        )
    )

    command.overwrite_output().run()

    print("Done!")

if __name__ == "__main__":
    main()
