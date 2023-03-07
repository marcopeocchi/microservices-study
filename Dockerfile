FROM alpine:3.17 AS build
WORKDIR /usr/src/fuu
RUN apk update && \
    apk add go nodejs npm gcc sqlite musl-dev
COPY . .
WORKDIR /usr/src/fuu/cmd/server/frontend
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

WORKDIR /etc/fuu
VOLUME /etc/fuu

RUN apk update && \
    apk add sqlite imagemagick ffmpeg

RUN chmod +x /usr/bin/fuu

ENV JWTSECRET=mykingmyknight
EXPOSE 4456
CMD exec /usr/bin/fuu
