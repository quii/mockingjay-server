FROM alpine:3.9.4

RUN apk add --update wget bash ca-certificates && \
    wget -O mockingjay-server https://github.com/quii/mockingjay-server/releases/download/1.11.1/linux_amd64_mockingjay-server --no-check-certificate && \
    chmod +x mockingjay-server && \
    apk del wget bash && \
    rm -rf /var/cache/apk/*

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENTRYPOINT ["./mockingjay-server"]
