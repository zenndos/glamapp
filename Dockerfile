# Use the official Golang image with version 1.20.4
FROM golang:1.20.4-alpine as build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the entire project directory into the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main ./src

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Expose the application port (replace with your actual port)
EXPOSE ${APP_PORT}

# Command to run the executable
CMD ["./main"]
