FROM golang:1.18.2-alpine3.16 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN mkdir /dist \
    && go build -o /dist/acmes


FROM alpine:3.16

COPY --from=builder /dist /

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone && mkdir /cert && mkdir /data

WORKDIR /

ENV ACMES_PORT 443
ENV ACMES_CA /cert/cert.pem
ENV ACMES_CAKEY /cert/key.pem
ENV ACMES_LOG_LEVEL info
ENV ACMES_STORE /data
ENV ACMES_EMAIL foo@acmes.com
ENV ACMES_DNS_PROVIDER alidns

VOLUME ["/cert", "/data"]

EXPOSE 443

ENTRYPOINT ["./acme", "serve"]

