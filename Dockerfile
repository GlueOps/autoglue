FROM golang:1.24-alpine@sha256:daae04ebad0c21149979cd8e9db38f565ecefd8547cf4a591240dc1972cf1399 AS builder

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