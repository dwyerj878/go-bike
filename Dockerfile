# --- Stage 1: Build ---
# Use an official Go image as the builder.
# Choose a specific Go version (e.g., 1.21) and Alpine for a smaller base.
FROM golang:1.24-alpine AS builder


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY src/ ./
RUN mkdir /static
RUN mkdir /test_data 
RUN mkdir /test_data/rider
RUN mkdir /test_data/data

COPY static/* /static/
COPY test_data/data/* /test_data/data/
COPY test_data/rider/* /test_data/rider/


# Build the Go app.
# - CGO_ENABLED=0: Build without Cgo for a static binary (usually preferred for containers)
# -o /app/main: Output the executable to /app/main
# ./... : Build all packages in the current directory and subdirectories (adjust if your main package is elsewhere)

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/main 

# --- Stage 2: Run ---
# Start from a minimal base image like Alpine.
# Using alpine is generally safer than scratch as it includes CA certificates and a shell.
# Use 'scratch' for the absolute smallest image if you are sure you don't need anything from an OS.
FROM alpine:latest
# FROM scratch

# Set the Current Working Directory inside the container
WORKDIR /app

# (Optional) Install CA certificates if your application makes HTTPS requests.
# Not needed if using 'scratch' base or if your app doesn't need certs.
# RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .


# (Optional) Copy static assets or templates if your application needs them
RUN mkdir /test_data
RUN mkdir /static
RUN mkdir /test_data/rider
RUN mkdir /test_data/data

COPY --from=builder /static/* /static
COPY --from=builder /test_data/data/* /test_data/data
COPY --from=builder /test_data/rider/* /test_data/rider

RUN ls -la /app
RUN ls -la /static
RUN ls -la /test_data


# Expose port 8080 to the outside world (Adjust this if your application listens on a different port)
EXPOSE 8081

# Command to run the executable
# Use ENTRYPOINT to make the container behave like an executable
#ENTRYPOINT ["./main"]
# Alternatively, use CMD if you want the command to be easily overridden:
CMD ["./main", "/test_data/data/activity_20250202.fit","/test_data/rider/rider.json","/test_data/data" ]
