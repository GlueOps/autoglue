#################################
# Builder: Go + Node in one
#################################
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache \
      git ca-certificates tzdata \
      build-base \
      nodejs npm

RUN npm i -g yarn pnpm

WORKDIR /src

COPY . .
RUN make swagger && make build


#################################
# Runtime
#################################
FROM alpine:3.21@sha256:b6a6be0ff92ab6db8acd94f5d1b7a6c2f0f5d10ce3c24af348d333ac6da80685

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