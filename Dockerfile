FROM alpine:3.14

MAINTAINER lichenxin

RUN apk --no-cache add ca-certificates curl bash xz-libs git
WORKDIR /tmp
RUN curl -L -O https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN tar -xf ffmpeg-release-amd64-static.tar.xz && \
      cd ff* && mv ff* /usr/local/bin && \
      rm -rf ffmpeg-release-amd64-static.tar.xz && \
      rm -rf ffmpeg-6.0-amd64-static

COPY ./bin /usr/local/bin/

EXPOSE 5050

CMD ["m3u8-http-server"]