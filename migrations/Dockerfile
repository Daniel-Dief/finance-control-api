FROM golang:1.27-alpine AS builder
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20
COPY --from=builder /go/bin/goose /usr/local/bin/goose
ENTRYPOINT ["goose"]
