FROM debian:stable-slim

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    git \
    openssh-client \
    python3 \
    python3-pip \
    python3-venv \
    rsync \
    ansible \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /root/.ssh && chmod 700 /root/.ssh

WORKDIR /workspace

ENTRYPOINT ["/bin/bash"]
