FROM golang:1.23-bullseye

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        make \
        python3 \
        python3-pip \
        ffmpeg \
    && rm -rf /var/lib/apt/lists/*

RUN pip3 install --no-cache-dir --upgrade pip && \
    pip3 install --no-cache-dir openai-whisper ffmpeg-python

WORKDIR /app

CMD ["/bin/bash"]
