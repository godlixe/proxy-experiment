# Use the official Go image as the base
FROM golang:1.23.5-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Use a minimal image for the final container
FROM debian:bullseye-slim

# Set the working directory inside the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the app
CMD ["./main"]
