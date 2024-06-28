# Start from the official Go image
FROM golang:1.21.5 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and /src files
COPY . .

# Change directory to src
WORKDIR /app/src

# Download all dependencies
RUN go mod download


# Build the Go application
RUN go build -o weather-monitoring .

# Expose port 8000 to the outside world
EXPOSE 8000

# Run the executable
CMD ["./weather-monitoring"]
