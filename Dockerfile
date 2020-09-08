FROM golang:1.15-alpine as builder

RUN apk add --no-cache gcc

ADD . /node

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN cd /node && go build -o myNode

FROM alpine:latest

COPY --from=builder /node/myNode /usr/local/bin/

ENTRYPOINT ["myNode"]
