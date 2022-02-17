# syntax=docker/dockerfile:1
FROM eqalpha/keydb:latest

RUN apt-get update -y && \
    apt-get install -y curl vim iproute2 less && \
    apt-get clean all

CMD keydb-server /etc/keydb/redis.conf --server-threads 2