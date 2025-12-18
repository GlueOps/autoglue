#################################
# Builder: Go + Node in one
#################################
FROM golang:1.25.5-alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder

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
FROM alpine:3.23@sha256:be171b562d67532ea8b3c9d1fc0904288818bb36fc8359f954a7b7f1f9130fb2

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