# Use an official Go image as a base
FROM golang:1.23.0

# Set the working directory inside the container
WORKDIR /config

# Copy go.mod and go.sum to download dependencies
COPY go.sum go.mod ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

RUN go build -o config ./cmd/main.go

# Expose the required port
EXPOSE 8080

# Run the service
CMD ["./config"]
