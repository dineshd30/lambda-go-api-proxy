# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go Lambda binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /main main.go

# Stage 2: Create the Lambda deployment image
FROM public.ecr.aws/lambda/go:1

# Copy the binary from the builder stage
COPY --from=builder /main ${LAMBDA_TASK_ROOT}

# Command to run the Lambda function
CMD ["main"]
