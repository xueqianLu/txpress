FROM golang:1.21-alpine AS build

ENV PACKAGES build-base

RUN apk add --update $PACKAGES

# Add source files
COPY ./ /build

RUN cd /build && go build -o /build/txpress

FROM alpine

WORKDIR /root

COPY  --from=build /build/txpress /usr/bin/txpress
COPY ./app.json /root/app.json
COPY ./accounts.json /root/accounts.json
ENTRYPOINT [ "/usr/bin/txpress --start --log /root/log/press.log " ]

