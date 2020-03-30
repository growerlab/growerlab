# growerlab services
#
FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y git supervisor

COPY ./data/supervisor/default.conf /etc/supervisor/supervisord.conf

VOLUME /data

WORKDIR /data

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
