# =========================
# 1) Build frontend (Vite)
# =========================
FROM node:20-alpine AS web-builder

WORKDIR /smart-contract
# Copy smart contracts (for ABI serving)
COPY smart-contract/abi ./abi

WORKDIR /web

# Install deps
COPY web/package.json web/package-lock.json ./
RUN npm install

# Copy frontend source and build SPA
COPY web .
RUN rm -f .env && npm run build


# =========================
# 2) Build backend (Go)
# =========================
# Use --platform=$BUILDPLATFORM to run Go compiler natively on the host architecture
FROM --platform=$BUILDPLATFORM golang:1.25.0-alpine AS go-builder

WORKDIR /app

# Cache Go deps
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Copy smart contracts (for ABI serving)
COPY smart-contract/abi ./smart-contract/abi

# Use Docker Buildx target variables ($TARGETOS, $TARGETARCH) for automatic multi-arch cross-compilation
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o server ./cmd/server


# =========================
# 3) Runtime image
# =========================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Frontend build output
COPY --from=web-builder /web/dist /app/web/dist

# Backend binary
COPY --from=go-builder /app/server /app/server

# Smart contract artifacts
COPY --from=go-builder /app/smart-contract/abi /app/smart-contract/abi

# Runtime envs
ENV GIN_MODE=release

EXPOSE 5001

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]