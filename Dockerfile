FROM golang:1.7-alpine
MAINTAINER Mainflux

ENV NATS_HOST nats
ENV NATS_PORT 4222

RUN apk update && apk add git && apk add wget && rm -rf /var/cache/apk/*

# copy the local package files into the container's workspace
COPY . /go/src/github.com/mainflux/mainflux-auth

# build the service inside the container
RUN go install github.com/mainflux/mainflux-auth

# Dockerize
ENV DOCKERIZE_VERSION v0.2.0
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

###
# Run main command with dockerize
###
CMD dockerize -wait tcp://$NATS_HOST:$NATS_PORT -timeout 10s /go/bin/mainflux-auth

# document exposed ports
EXPOSE 8180
