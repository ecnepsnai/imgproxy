FROM docker.io/library/alpine:3
LABEL maintainer="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.authors="Ian Spence <ian@ecnepsnai.com>"
LABEL org.opencontainers.image.source=https://github.com/ecnepsnai/imgproxy
LABEL org.opencontainers.image.title="imgproxy"
LABEL org.opencontainers.image.description="HTTP Image Proxy"

EXPOSE 80 443

ADD imgproxy /imgproxy
ENTRYPOINT [ "/imgproxy" ]
