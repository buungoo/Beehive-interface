# Start from the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /

# Copy go.mod and go.sum first to cache dependency installation
COPY go.mod ./

# Clean up unused dependancy before download
#RUN go mod tidy

# Download dependencies (caches the downloaded modules)
RUN go mod download

# Copy the Go project to the container
COPY . .



# Build the Go API
RUN go build -o main ./main.go

# Expose port 8080 (adjust if your Go app uses a different port)
EXPOSE 8080

# Command to run the Go API
CMD ["./main"]
