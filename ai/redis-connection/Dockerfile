# Multi-stage Dockerfile pinned to Go 1.24.5
# Build stage
FROM golang:1.24.5 AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o /out/redis-conn .

# Run stage (distroless/static base)
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/redis-conn /redis-conn
USER nonroot:nonroot
ENTRYPOINT ["/redis-conn"]
