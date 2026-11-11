# ==============================================================================
# Stage 1: Build React 19 + TypeScript + Tailwind v4 Frontend Bundle
# ==============================================================================
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# ==============================================================================
# Stage 2: Build Pure-Go Backend Binary (CGO_ENABLED=0)
# ==============================================================================
FROM golang:1.27-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /eneba-server ./cmd/server

# ==============================================================================
# Stage 3: Minimal Hardened Production Runtime (Non-Root User + Secret Injection)
# ==============================================================================
FROM alpine:3.21 AS runtime
RUN apk --no-cache add ca-certificates tzdata wget && \
    addgroup -S eneba && adduser -S eneba -G eneba

WORKDIR /app
RUN mkdir -p /app/data /app/static && chown -R eneba:eneba /app

COPY --from=backend-builder --chown=eneba:eneba /eneba-server /app/eneba-server
COPY --from=frontend-builder --chown=eneba:eneba /app/frontend/dist /app/static

USER eneba:eneba

ENV APP_ENV=production \
    PORT=8080 \
    DB_PATH=/app/data/eneba.db \
    STATIC_DIR=/app/static

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/api/v1/health || exit 1

ENTRYPOINT ["/app/eneba-server"]
