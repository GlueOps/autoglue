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

RUN npm i -g yarn pnpm

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
COPY --from=builder /src/autoglue /app/autoglue

ENV PORT=8080
EXPOSE 8080
USER app

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT}/api/v1/healthz" || exit 1

ENTRYPOINT ["/app/autoglue"]