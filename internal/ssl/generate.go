package ssl

import (
	"fmt"
	"github.com/aacfactory/afssl"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generate(cn string, expires int, out string) (err error) {
	if cn == "" {
		cn = "acmes"
	}
	if expires < 1 {
		expires = 365
	}
	if out == "" {
		out = "."
	}
	out, err = filepath.Abs(out)
	if err != nil {
		err = fmt.Errorf("acmes: generate ca failed, %v", err)
		return
	}
	outExist := false
	_, outStatErr := os.Stat(out)
	if outStatErr == nil {
		outExist = true
	} else {
		outExist = !os.IsNotExist(outStatErr)
	}
	if !outExist {
		mkdirErr := os.MkdirAll(out, 0600)
		if mkdirErr != nil {
			err = fmt.Errorf("acmes: generate ca failed for create out dir failed, %v", mkdirErr)
			return
		}
	}
	config := afssl.CertificateConfig{
		Country:            "",
		Province:           "",
		City:               "",
		Organization:       "",
		OrganizationalUnit: "",
		CommonName:         cn,
		IPs:                nil,
		Emails:             nil,
		DNSNames:           nil,
	}
	caPEM, caKeyPEM, caErr := afssl.GenerateCertificate(config, afssl.CA(), afssl.WithExpirationDays(expires))
	if caErr != nil {
		err = fmt.Errorf("acmes: generate ca failed, %v", caErr)
		return
	}
	certPath := filepath.Join(out, "cert.pem")
	saveCertErr := ioutil.WriteFile(certPath, caPEM, 0600)
	if saveCertErr != nil {
		err = fmt.Errorf("acmes: generate ca failed, %v", saveCertErr)
		return
	}
	keyPath := filepath.Join(out, "key.pem")
	saveKeyErr := ioutil.WriteFile(keyPath, caKeyPEM, 0600)
	if saveKeyErr != nil {
		err = fmt.Errorf("acmes: generate ca failed, %v", saveKeyErr)
		return
	}
	return
}
