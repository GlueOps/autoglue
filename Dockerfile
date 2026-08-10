#################################
# Builder: Go + Node in one
#################################
FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder

RUN apk add --no-cache \
    bash git ca-certificates tzdata \
    build-base \
    nodejs npm \
    openjdk17-jre-headless \
    jq yq brotli

# The UI is on Yarn 4 (berry). `npm i -g yarn` would install classic 1.x, which
# cannot read a berry lockfile and ignores ui/.yarnrc.yml entirely. Yarn 4 is
# not published to npm either, so corepack is the only route — and Alpine's
# nodejs package does not bundle it.
#
# The exact version comes from "packageManager" in ui/package.json.
RUN npm i -g corepack@0.35.0 && corepack enable

# Corepack prompts before downloading a package manager; that would hang a build.
ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0

WORKDIR /src

COPY . .
RUN make clean
RUN make swagger
RUN make sdk-ts-ui
RUN make ui
RUN make build

#################################
# Runtime
#################################
FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

RUN apk add --no-cache ca-certificates tzdata postgresql17-client \
 && addgroup -S app && adduser -S app -G app

WORKDIR /app

# Install onto PATH, not just /app. Kubernetes `command:` replaces ENTRYPOINT
# rather than appending to it, so a manifest saying `command: ["autoglue",
# "worker"]` resolves the name through $PATH — and /app is not on it:
#   /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
# Without this the container dies with
#   exec: "autoglue": executable file not found in $PATH
COPY --from=builder /src/autoglue /usr/local/bin/autoglue

# Keep the historical path working for anything that hardcodes /app/autoglue.
RUN ln -s /usr/local/bin/autoglue /app/autoglue

ENV PORT=8080
EXPOSE 8080
USER app

# Only meaningful for the `serve` command. A container started with `worker`
# serves no HTTP and should override or disable this healthcheck.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT}/api/v1/healthz" || exit 1

# One image, two roles. Override the command with `worker` to run background
# jobs; `serve` never processes them.
#
# Works as `args: ["worker"]` (ENTRYPOINT is kept) and as
# `command: ["autoglue", "worker"]` (ENTRYPOINT is replaced, name found on PATH).
ENTRYPOINT ["autoglue"]
CMD ["serve"]
