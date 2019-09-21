FROM golang:1.13.0-alpine AS builder

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build .

# Create the minimal runtime image
FROM alpine:3.10

EXPOSE 8080
COPY --chown=0:0 --from=builder /build /app

# Set up the app to run as a non-root user inside the /data folder
# User ID 65534 is usually user 'nobody'.
# The executor of this image should still specify a user during setup.
RUN chown -R 65534 /app
USER 65534
WORKDIR /app

CMD "./Server"
