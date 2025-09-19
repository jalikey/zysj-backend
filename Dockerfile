# --- Stage 1: Build ---
# Use an official Go image as the build environment.
FROM golang:1.25.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./
#RUN export GOPROXY=https://goproxy.cn,direct
#RUN export all_proxy="socks5://127.0.0.1:10808"
RUN export GOPROXY=https://goproxy.cn,direct && go mod download
#RUN unset all_proxy
# Copy the entire source code
COPY . .

# Build the application, disabling CGO for a static binary.
# The output binary will be named 'server'.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/api

# --- Stage 2: Final Image ---
# Use a minimal base image for the final container.
# alpine is a good choice for its small size.
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/server .

# We don't need any other source files.
# The app will read config from environment variables, not a .env file.

# Expose the port the app runs on.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["./server"]