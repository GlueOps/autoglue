#################################
# Builder: Go + Node in one
#################################
FROM golang:1.25-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS builder

RUN apk add --no-cache \
      git ca-certificates tzdata \
      build-base \
      nodejs npm

RUN npm i -g yarn pnpm

WORKDIR /src

COPY . .
RUN make clean && make swagger && make ui && make build


#################################
# Runtime
#################################
FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /src/autoglue /app/autoglue

ENV PORT=8080
EXPOSE 8080
USER app

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT}/api/healthz" || exit 1

ENTRYPOINT ["/app/autoglue"]