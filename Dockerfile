# growerlab services
#
FROM ubuntu:latest

RUN apt-get update && apt-get install -y supervisor

COPY ./data/supervisor/default.conf /etc/supervisor/supervisord.conf

VOLUME /data
VOLUME /data/logs

WORKDIR /data

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
