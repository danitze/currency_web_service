FROM golang:1.22.3

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o /currency-web-service

# Expose port 8080
EXPOSE 8080

# Set the default command to run when the container starts
CMD ["/currency-web-service"]
