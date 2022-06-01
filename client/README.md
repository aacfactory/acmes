# Acmes Client
## Install
```shell
go get github.com/aacfactory/acmes/client
```
## Usage
```go
ca, _ := ioutil.ReadFile("./cert.pem")
key, _ := ioutil.ReadFile("./key.pem")

acme, err := client.New(ca, key, "127.0.0.1:8443")
if err != nil {
    t.Error(err)
    return
}
_, cancel, obtainErr := acme.Obtain(context.TODO(), "*.foo.com")
if obtainErr != nil {
    t.Error(obtainErr)
    return
}
// to cancel auto renew
cancel()
```