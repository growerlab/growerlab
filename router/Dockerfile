# growerlab router for dev
#
FROM scratch

COPY ./main /home/router

RUN chmod +x /home/router

VOLUME /data

WORKDIR /home

ENTRYPOINT ["/home/router"]
