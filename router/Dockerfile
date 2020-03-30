# growerlab services
#
FROM scratch

RUN go build -o /tmp/router router.go

COPY /tmp/router /home/router

RUN chmod +x /home/router

VOLUME /data

WORKDIR /home

ENTRYPOINT ["/home/router"]
