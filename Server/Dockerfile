FROM golang:1.13.0-alpine AS builder

RUN apk add build-base # for gcc

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build cmd/http/main.go

# Create the minimal runtime image
FROM alpine:3.10

WORKDIR /app

# Bring the binary over
COPY --from=builder /build/main /app

# Set up the app to run as a non-root user
RUN chown -R 65534 /app
USER 65534

EXPOSE 8080
CMD "./main"
