#################################
# Builder: Go + Node in one
#################################
FROM golang:1.25.4-alpine@sha256:d3f0cf7723f3429e3f9ed846243970b20a2de7bae6a5b66fc5914e228d831bbb AS builder

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
FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

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