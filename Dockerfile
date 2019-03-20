FROM golang:alpine as build
RUN apk --no-cache --update upgrade && apk --no-cache add git

ADD . /go/src/github.com/mback2k/simple-file-server
WORKDIR /go/src/github.com/mback2k/simple-file-server

RUN go get
RUN go build -ldflags="-s -w"
RUN chmod +x simple-file-server

FROM mback2k/alpine:latest
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates

COPY --from=build /go/src/github.com/mback2k/simple-file-server/simple-file-server /usr/local/bin/simple-file-server

RUN addgroup -S serve
RUN adduser -h /data -S -D -G serve serve

WORKDIR /data
USER serve

CMD [ "/usr/local/bin/simple-file-server" ]