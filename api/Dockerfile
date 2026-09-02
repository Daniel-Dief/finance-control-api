FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/finance-control-api .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/finance-control-api .

EXPOSE 1323

USER nonroot:nonroot

ENTRYPOINT ["./finance-control-api"]
