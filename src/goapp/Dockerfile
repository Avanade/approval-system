FROM golang:1.23.2-bookworm

# Set the Current Working Directory inside the container
WORKDIR /goapp

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

EXPOSE 8080

# Build the Go app
RUN go build -o approvalsystem .

# Command to run the executable
CMD ./approvalsystem