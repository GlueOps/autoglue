FROM golang:1.24-alpine AS builder

RUN apk update && apk add make git

WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN make swagger && go build -o autoglue main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/autoglue /app/autoglue

ENV PORT=8080

EXPOSE 8080

CMD ["./autoglue", "serve"]