# Use the official Go image
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy the Go modules and app source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Expose the app port
EXPOSE 8080

# Run the Go application
CMD ["go", "run", "main.go"]
