FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

# Automatic platform args set by buildx (e.g. --platform linux/amd64,linux/arm64).
# Don't assign defaults here: that would override buildx's automatic values.
ARG TARGETOS
ARG TARGETARCH

# go install writes cross-compiled binaries to $GOPATH/bin/$GOOS_$GOARCH, but
# native ones straight to $GOPATH/bin, so normalize the output path.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go install github.com/pressly/goose/v3/cmd/goose@latest \
 && mkdir -p /out \
 && cp "$(find /go/bin -type f -name goose | head -n 1)" /out/goose

FROM alpine:3.20
COPY --from=builder /out/goose /usr/local/bin/goose
ENTRYPOINT ["goose"]
