FROM alpine:3.17 AS build
WORKDIR /usr/src/fuu
RUN apk update && \
    apk add go nodejs npm gcc sqlite musl-dev
COPY . .
WORKDIR /usr/src/fuu/frontend
RUN npm i
RUN npm run build
WORKDIR /usr/src/fuu
RUN CGO_ENABLED=1 go build -o fuu main.go


FROM alpine:3.17
COPY --from=build /usr/src/fuu /usr/bin

WORKDIR /media
VOLUME /media
WORKDIR /cache
VOLUME /cache

RUN apk update && \
    apk add sqlite imagemagick ffmpeg

ENV MASTERPASS=adminadminadmin
ENV SECRET=secret
ENV THUMBNAIL_HEIGHT=450
ENV THUMBNAIL_QUALITY=75

RUN chmod +x /usr/bin/fuu

EXPOSE 4456
CMD exec /usr/bin/fuu -w /media -S $SECRET -M $MASTERPASS
