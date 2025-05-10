FROM golang:1-bookworm AS dev

ENV DEBIAN_FRONTEND=noninteractive

# Update the package list and install FFmpeg
RUN apt-get update && \
  apt-get install -y ffmpeg && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

CMD ["bash"]
