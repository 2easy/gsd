# Build stage
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/ ./
RUN npm ci && npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
# Add build dependencies
RUN apk add --no-cache gcc musl-dev
COPY . .
COPY --from=frontend-builder /app/dist ./backend/dist
RUN CGO_ENABLED=1 go build -ldflags='-linkmode external -extldflags "-static"' -o gsd ./backend

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/gsd .
VOLUME /data
EXPOSE 8080

# Set default arguments that can be overridden at runtime
ENTRYPOINT ["/app/gsd", "--port", "8080", "--db", "/data/gsd.db"]