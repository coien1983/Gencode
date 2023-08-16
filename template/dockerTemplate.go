package template

var DockerTemplate = `FROM golang:1.18.7-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY . .

RUN go build -o %s main.go

EXPOSE %s

CMD ["./%s"]`
