# Stage 1: Build SvelteKit frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.23-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
# Enable CGO for go-sqlite3
ENV CGO_ENABLED=1
RUN go build -o /app/wg-manager ./cmd/server/main.go

# Stage 3: Final lightweight image
FROM alpine:latest
RUN apk add --no-cache wireguard-tools iptables
WORKDIR /app

# Copy binaries and built static files
COPY --from=backend-builder /app/wg-manager /app/wg-manager
COPY --from=frontend-builder /app/build /app/public

# Copy default config template if needed (optional)
COPY backend/internal/config/config.json /app/internal/config/config.json

# Expose HTTP port and WireGuard port
EXPOSE 8080
EXPOSE 51820/udp

# Entrypoint for the Go server
CMD ["/app/wg-manager"]
