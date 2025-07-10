# Build the binary
FROM golang:1.18.3 as builder
WORKDIR /go/src/github.com/my-org/pwk
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
RUN go build -o mutating-admission-webhook ./main.go

# Copy the binary into a thin image
FROM alpine:3.13.5
RUN apk add gcompat
WORKDIR /root
COPY server.pem server.pem
COPY server-key.pem server-key.pem
COPY --from=builder /go/src/github.com/my-org/pwk/mutating-admission-webhook /usr/local/bin/mutating-admission-webhook