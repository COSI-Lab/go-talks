# Build the server binary
FROM golang:latest AS go-build

WORKDIR /go/src/

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -v -o server

# Copies the server binary into a minimal image
FROM debian:latest

WORKDIR /root/

# Add config, templates, and static files
COPY ./config.toml .
COPY ./templates ./templates
COPY ./static ./static
COPY ./posts ./posts

COPY --from=go-build /go/src/server .

# Run the server
CMD ["/root/server"]