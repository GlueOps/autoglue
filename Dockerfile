FROM golang:1.24-alpine@sha256:c8c5f95d64aa79b6547f3b626eb84b16a7ce18a139e3e9ca19a8c078b85ba80d AS builder

RUN apk update && apk add make git

WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN make swagger && go build -o autoglue main.go

FROM alpine:latest@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

WORKDIR /app
COPY --from=builder /app/autoglue /app/autoglue

ENV PORT=8080

EXPOSE 8080

CMD ["./autoglue", "serve"]