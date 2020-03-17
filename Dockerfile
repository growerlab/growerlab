# growerlab services
#
FROM ubuntu:latest

RUN apt-get update && apt-get install -y supervisor
COPY ./data/supervisor/*.conf /etc/supervisor/conf.d/

VOLUME /data

WORKDIR /data

CMD ["/usr/bin/supervisord","-c","/etc/supervisor/supervisord.conf"]
