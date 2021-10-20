FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod go.sum main.go ./
COPY handler/ ./handler/

RUN go mod download
RUN go build -o /dcbot -ldflags '-s -w'


FROM jrottenberg/ffmpeg:4.1-alpine as runner

COPY --from=builder /dcbot /dcbot

ENTRYPOINT ["/dcbot"]
