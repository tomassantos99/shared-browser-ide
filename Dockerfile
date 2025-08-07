# Stage 1: Build React frontend
FROM node:20-alpine AS frontend

WORKDIR /app
COPY web/shared-browser-ide ./shared-browser-ide
WORKDIR /app/shared-browser-ide

RUN npm install
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24.5-alpine AS backend

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy React build to be embedded/served by Go
COPY --from=frontend /app/shared-browser-ide/dist ./web/shared-browser-ide/dist

# Build Go binary
WORKDIR /app/cmd/server
RUN go build -o app

# Final Stage: Minimal container to run the app
FROM alpine:latest

WORKDIR /root/

COPY --from=backend /app/cmd/server/app .
COPY --from=backend /app/web/shared-browser-ide/dist ./static

EXPOSE 8080

CMD ["./app"]
