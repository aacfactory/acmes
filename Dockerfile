FROM golang:1.21.6-alpine3.19 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN mkdir /dist \
    && go build -o /dist/acmes


FROM alpine:3.19

COPY --from=builder /dist /

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone && mkdir /cert && mkdir /data && chmod +x /acmes

WORKDIR /

ENV ACMES_PORT 443
ENV ACMES_CA /cert/cert.pem
ENV ACMES_CAKEY /cert/key.pem
ENV ACMES_LOG_LEVEL debug
ENV ACMES_STORE file:///data
ENV ACMES_EMAIL foo@acmes.com
ENV ACMES_DNS_PROVIDER alidns

VOLUME ["/cert", "/data"]

EXPOSE 443

ENTRYPOINT ["./acmes", "serve"]

