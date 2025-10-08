# builder
FROM golang:1.25 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/api/

# runtime
FROM gcr.io/distroless/static
COPY --from=builder /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
