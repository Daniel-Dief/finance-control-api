FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

# Automatic platform args set by buildx (e.g. --platform linux/amd64,linux/arm64).
# No default here: a default value would override the built-in automatic ARG.
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags="-s -w" -o /app/finance-control-api .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/finance-control-api .

EXPOSE 1323

USER nonroot:nonroot

ENTRYPOINT ["./finance-control-api"]

