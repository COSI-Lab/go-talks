# build web server
FROM golang:latest AS go-build

WORKDIR /go/src/

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./*.go ./
RUN go build -o server

# copy server and static files to clean alpine image
FROM debian:latest

WORKDIR /srv/website
EXPOSE 5000

RUN apt-get update -y && apt-get upgrade -y && apt-get dist-upgrade -y && apt-get install curl -y
RUN curl -sSL https://get.docker.com/ | sh

COPY --from=go-build /go/src/ ./
COPY ./static ./static


CMD ["./server"]