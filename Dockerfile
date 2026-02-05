#################################
# Builder: Go + Node in one
#################################
FROM golang:1.25.7-alpine@sha256:81d49e1de26fa223b9ae0b4d5a4065ff8176a7d80aa5ef0bd9f2eee430afe4d7 AS builder

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
FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62

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