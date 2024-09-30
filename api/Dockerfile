# Start from the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /

# Copy the Go project to the container
COPY . .

# Download dependencies
# RUN go mod tidy

# Build the Go API
RUN go build -o main ./main.go

# Expose port 8080 (adjust if your Go app uses a different port)
EXPOSE 8080

# Command to run the Go API
CMD ["./main"]
