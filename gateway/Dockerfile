FROM golang:1.20-alpine

ENV GOPATH=/

COPY ./ ./


RUN go mod download
RUN go build -o gateway ./cmd/main.go


CMD ["./gateway"]