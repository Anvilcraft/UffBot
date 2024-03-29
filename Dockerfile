FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

# Copy the code into the container
COPY . .

# Run test
RUN go test ./...

# Build the application
RUN go build main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .
RUN cp /build/config.json .

CMD ["/dist/main","-t","USER-TOKEN"]
