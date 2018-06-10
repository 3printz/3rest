FROM golang:1.9

MAINTAINER Eranga Bandara (erangaeb@gmail.com)

# install dependencies
RUN go get github.com/gocql/gocql
RUN	go get github.com/gorilla/mux
RUN go get github.com/Shopify/sarama
RUN go get github.com/wvanbergen/kafka/consumergroup

# env
ENV SENZIE_NAME sampath
ENV SENZIE_MODE DEV
ENV CASSANDRA_HOST dev.localhost
ENV CASSANDRA_PORT 9042
ENV CASSANDRA_KEYSPACE zchain
ENV KAFKA_TOPIC sampath
ENV KAFKA_CGROUP sampathg
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
