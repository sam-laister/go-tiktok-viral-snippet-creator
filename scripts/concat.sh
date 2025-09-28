#!/usr/bin/env zsh
set -euo pipefail

# Directory with mp4 files (default: current directory)
DIR="${1:-.}"

# Output directory
OUT="$DIR/concatenated"
mkdir -p "$OUT"

# Gather and sort .mp4 files into an array
FILES=("${(@f)$(find "$DIR" -maxdepth 1 -type f -name '*.mp4' | sort)}")

# Loop through files in pairs
i=1
while (( i <= ${#FILES} - 1 )); do
    f1="${FILES[$i]}"
    f2="${FILES[$((i+1))]}"
    base1="${${f1:t}%.*}"
    base2="${${f2:t}%.*}"
    out_file="$OUT/${base1}_${base2}.mp4"

    echo "Concatenating: $f1 + $f2 -> $out_file"

    # Create a temporary concat list file
    tmpfile=$(mktemp /tmp/concatlist.XXXXXX)
    print "file '$f1'\nfile '$f2'" > "$tmpfile"

    # Concatenate without re-encoding (requires same codec/parameters)
    ffmpeg -hide_banner -loglevel error -f concat -safe 0 -i "$tmpfile" -c copy "$out_file"

    rm "$tmpfile"
    (( i += 2 ))
done

