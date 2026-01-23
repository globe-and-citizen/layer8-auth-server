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

# Copy frontend source
COPY web .
RUN rm .env

# Build SPA
RUN npm run build


# =========================
# 2) Build backend (Go)
# =========================
FROM golang:1.25.0-alpine AS go-builder

WORKDIR /app

#RUN apk add --no-cache git

# Cache Go deps
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# Copy smart contracts (for ABI serving)
COPY smart-contract/abi ./smart-contract/abi

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

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

# Runtime envs (read by main.go)
ENV GIN_MODE=release
#ENV PORT=5001

EXPOSE 5001

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
