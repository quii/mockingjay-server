FROM alpine:3.3
RUN apk add --update wget bash
RUN wget -O mockingjay-server https://github.com/quii/mockingjay-server/releases/download/1.5.2/linux_amd64_mockingjay-server --no-check-certificate
RUN chmod +x mockingjay-server
ENTRYPOINT ["./mockingjay-server"]
