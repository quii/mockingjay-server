FROM progrium/busybox
RUN opkg-install wget bash
RUN wget -O mockingjay-server https://github.com/quii/mockingjay-server/releases/download/1.5.2/linux_amd64_mockingjay-server --no-check-certificate
RUN chmod +x mockingjay-server
ENTRYPOINT ["./mockingjay-server"]
