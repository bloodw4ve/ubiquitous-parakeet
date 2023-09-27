FROM golang:1.20-alpine

ENV GOPATH=/

COPY ./ ./


RUN go mod download
RUN go build -o news ./cmd/main.go


CMD ["./news"]