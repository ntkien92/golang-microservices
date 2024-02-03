FROM golang:1.21-alpine as build

WORKDIR $GOPATH/github.com/ntkien92/golang-microservices

RUN apk update

RUN apk add --no-cache gcc musl-dev linux-headers git

ADD go.mod .
ADD go.sum .

RUN go mod download

ADD . .

RUN go install github.com/ntkien92/golang-microservices

# ---


FROM alpine:3.19
ARG dist=0.0
COPY --from=build /go/bin/golang-microservices /

ENV LOG_LEVEL=INFO
ENV GIN_MODE=release
ENV SERVER_VERSION=$dist

EXPOSE 8000

CMD ["/golang-microservices"]
