# Build stage
FROM public.ecr.aws/docker/library/golang:1.22-alpine3.20 as builder

# Install git and Delve
RUN apk add --no-cache git bash

# Create build directory and set working directory
RUN mkdir /build
WORKDIR /build

# Set up GITHUB_TOKEN for private repos if needed
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"

# Copy the project files
COPY . .

# Tidy the Go modules
RUN go mod tidy

# Install Delve for debugging
#RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Build the Go project
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ticketon-auth-service .

# Final stage (runtime)
FROM alpine

# Install libc6-compat for compatibility and bash
RUN apk add --no-cache build-base bash curl

# Copy the built Go binary and .env file
COPY --from=builder /ticketon-auth-service /ticketon-auth-service
COPY --from=builder /build/.env.example /.env

# Copy Delve debugger from the builder stage
#COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

# Expose port for Delve debugger
#EXPOSE 2345

# Start the app with Delve debugger for remote debugging
#ENTRYPOINT ["dlv", "exec", "/ticketon-auth-service", "--headless", "--listen=:2345", "--api-version=2", "--log"]
ENTRYPOINT ["./ticketon-auth-service"]
