FROM alpine:3.5

ADD myip /

ENTRYPOINT [ "/myip" ]

