# Acmes
acme server

## Install
```shell
go install github.com/aacfactory/acmes
```
## Usage
Setup ACME DNS Provider. More DNS providers is [HERE](https://go-acme.github.io/lego/dns/).
```shell
export ALICLOUD_ACCESS_KEY=foo
export ALICLOUD_SECRET_KEY=bar
```
Startup server.
```shell
acmes serve --port 8443 \
  --ca ./cert.pem --cakey ./key.pem \
  --level debug \
  --store file:///some_path/store \
  --provider alidns \
  --email for@bar.com 
```
Run in docker
```shell
docker run -d --rm --name acmes \
  -e ACMES_EMAIL=foo@bar.com \
  -e ACMES_DNS_PROVIDER=alidns \
  -e ALICLOUD_ACCESS_KEY=foo \
  -e ALICLOUD_SECRET_KEY=bar \
  -v $PWD/data=/data \
  -v $PWD/cert=/cert \
  wangminxiang0425/acmes
```
Use Client in your project, see `client`.