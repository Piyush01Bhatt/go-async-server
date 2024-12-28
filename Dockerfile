# Use Ubuntu as the base image
FROM golang:1.23.2

# Install required packages
RUN apt-get update

# Set the working directory
WORKDIR /app

# Copy the Go application into the container
COPY . .

# Build the Go application
RUN go build -o async-server cmd/server.go

# Command to run your application
CMD ["./async-server"]