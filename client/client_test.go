package client_test

import (
	"context"
	"github.com/aacfactory/acmes/client"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	ca, _ := os.ReadFile("G:\\test_acmes\\cert.pem")
	key, _ := os.ReadFile("G:\\test_acmes\\key.pem")

	acme, err := client.New(ca, key, "127.0.0.1:8443")
	if err != nil {
		t.Error(err)
		return
	}
	_, cancel, obtainErr := acme.Obtain(context.TODO(), "*.aacfactory.com")
	if obtainErr != nil {
		t.Error(obtainErr)
		return
	}
	cancel()
}
