FROM golang:1.7-alpine
MAINTAINER Mainflux

env NATS_HOST nats
env NATS_PORT 4222

env REDIS_HOST redis
env REDIS_PORT 6379

RUN apk update && apk add git && apk add wget && rm -rf /var/cache/apk/*

# copy the local package files into the container's workspace
COPY . /go/src/github.com/mainflux/mainflux-auth

RUN mkdir -p /etc/mainflux/auth
COPY config/config-docker.toml /etc/mainflux/auth/config.toml

# build the service inside the container
RUN go install github.com/mainflux/mainflux-auth

# Dockerize
ENV DOCKERIZE_VERSION v0.2.0
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

###
# Run main command with dockerize
###
CMD dockerize -wait tcp://$NATS_HOST:$NATS_PORT \
					-wait tcp://$REDIS_HOST:$REDIS_PORT \
					-timeout 10s /go/bin/mainflux-auth -c /etc/mainflux/auth/config.toml

# document exposed ports
EXPOSE 8180
