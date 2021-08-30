FROM golang:1.13.5 AS build

WORKDIR /low-code

ENV GOPROXY https://goproxy.io,direct
ENV GO111MODULE on

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN mkdir -p dist \
    && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode "external" -extldflags "-static"' -o ./dist/low-code ./cmd/main.go \
    && mkdir -p /app \
    && cp -r resource /app/ \
    && mv dist/* /app/


FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=build /app .
USER root
CMD ["./low-code"]
