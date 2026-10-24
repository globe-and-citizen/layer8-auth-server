# =========================
# 1) Build frontend (Vite)
# =========================
FROM node:20-alpine AS web-builder

WORKDIR /smart-contract
# Copy smart contracts (for ABI serving)
COPY smart-contract/abi ./abi

WORKDIR /web

# Install dependencies
COPY web/package.json web/package-lock.json ./
RUN npm install

# Copy frontend source and build SPA
COPY web .
RUN rm -f .env && npm run build


# =========================
# 2) Build backend (Go)
# =========================
# Run Go compiler natively on the build host
FROM --platform=$BUILDPLATFORM golang:1.25.11-alpine AS go-builder

WORKDIR /app

# Cache Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Copy smart contracts (for ABI serving)
COPY smart-contract/abi ./smart-contract/abi

# Buildx target variables for multi-arch compilation
ARG TARGETOS
ARG TARGETARCH

# Build application server
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o server ./cmd/server

# Build database migration executable
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o migrate ./cmd/migrate


# =========================
# 3) Runtime image
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Frontend build output
COPY --from=web-builder /web/dist /app/web/dist

# Backend binary
COPY --from=go-builder /app/server /app/server

# Migration binary
COPY --from=go-builder /app/migrate /app/migrate

# Migration SQL files
COPY --from=go-builder \
    /app/cmd/migrate/migrations \
    /app/migrations

# Smart contract artifacts
COPY --from=go-builder \
    /app/smart-contract/abi \
    /app/smart-contract/abi

# Runtime environment
ENV GIN_MODE=release

EXPOSE 5001

# Run as non-root
USER nonroot:nonroot

# Default command: application server
# ECS can override this to ["/app/migrate"]
CMD ["/app/server"]
