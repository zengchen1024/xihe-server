FROM golang:latest as BUILDER

# build binary
COPY . /go/src/github.com/opensourceways/xihe-server
RUN cd /go/src/github.com/opensourceways/xihe-server && go build -v -o ./xihe-server main.go

# copy binary config and utils
FROM alpine:latest
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/xihe-server/xihe-server /opt/app
COPY ./deploy/conf.yaml /opt/app/conf

ENTRYPOINT ["/opt/app/xihe-server"]
