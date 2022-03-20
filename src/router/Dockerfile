# growerlab router for dev
#
FROM ubuntu:latest

COPY ./main /home/router

RUN chmod +x /home/router

VOLUME /data

WORKDIR /home

ENTRYPOINT ["/home/router"]
