FROM golang:1.25-alpine@sha256:f18a072054848d87a8077455f0ac8a25886f2397f88bfdd222d6fafbb5bba440 AS builder

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