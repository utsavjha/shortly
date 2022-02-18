# syntax=docker/dockerfile:1
FROM eqalpha/keydb:latest

RUN apt-get update -y && \
    apt-get install -y curl vim iproute2 less && \
    apt-get clean all

COPY k2.sh /opt

CMD keydb-server /etc/keydb/redis.conf --server-threads 2