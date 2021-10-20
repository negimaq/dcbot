FROM golang:1.16-alpine

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY go.mod go.sum main.go ./
COPY handler/ ./handler/

RUN go mod download
RUN go build -o /dcbot

CMD ["/dcbot"]
