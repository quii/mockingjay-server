FROM alpine:3.3

RUN apk add --update wget bash && \
    wget -O mockingjay-server https://github.com/quii/mockingjay-server/releases/download/1.5.3/linux_amd64_mockingjay-server --no-check-certificate && \
    chmod +x mockingjay-server && \
    apk del wget bash && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["./mockingjay-server"]
