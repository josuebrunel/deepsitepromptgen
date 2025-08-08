# Build
FROM    golang:latest AS build

ENV     GO111MODULE=on
ENV     CGO_ENABLED=0

RUN     mkdir /go/src/app
WORKDIR /go/src/app
ADD     . /go/src/app
RUN     go get github.com/labstack/echo/v5@v5.0.0-20230722203903-ec5b858dab61
RUN     make deps
RUN     make build
EXPOSE  8080 4000


# Deploy
FROM    alpine:latest
RUN     mkdir /opt/dspg
WORKDIR /opt/dspg
COPY    --from=build /go/src/app/bin/dspg /opt/dspg/dspg
EXPOSE  8080
CMD     ["./dspg"]
