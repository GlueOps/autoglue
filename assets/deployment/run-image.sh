#!/usr/bin/env bash

docker run --rm -it \
  -v "$HOME/.ssh:/root/.ssh:ro" \
  -v "$(pwd)/ansible:/workspace" \
  ansible-debian:latest
