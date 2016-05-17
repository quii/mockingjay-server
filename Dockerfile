FROM alpine:3.3

RUN apk add --update wget bash ca-certificates && \
    wget -O mockingjay-server https://github.com/quii/mockingjay-server/releases/download/1.5.8/linux_amd64_mockingjay-server --no-check-certificate && \
    chmod +x mockingjay-server && \
    apk del wget bash && \
    rm -rf /var/cache/apk/*

ENTRYPOINT ["./mockingjay-server"]
