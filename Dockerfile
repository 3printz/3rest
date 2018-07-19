FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/gocql/gocql
RUN	go get github.com/gorilla/mux
RUN go get github.com/Shopify/sarama
RUN go get github.com/wvanbergen/kafka/consumergroup

# env
ENV SENZIE_NAME restz
ENV SENZIE_MODE DEV
ENV KAFKA_TOPIC restz
ENV KAFKA_CGROUP restzg
ENV KAFKA_KHOST dev.localhost
ENV KAFKA_KPORT 9092
ENV KAFKA_ZHOST dev.localhost
ENV KAFKA_KPORT 2181

# copy app
ADD . /app
WORKDIR /app

# build
RUN go build -o build/senz src/*.go

# server running port
EXPOSE 7070

# .keys volume
VOLUME ["/app/.keys"]

ENTRYPOINT ["/app/docker-entrypoint.sh"]
