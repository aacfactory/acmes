package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/aacfactory/afssl"
	"os"
)

func createTLSConfig(ca string, key string) (config *tls.Config, err error) {
	caPEM, caErr := os.ReadFile(ca)
	if caErr != nil {
		err = fmt.Errorf("acmes: read ca file failed, %v", caErr)
		return
	}
	keyPEM, keyErr := os.ReadFile(key)
	if keyErr != nil {
		err = fmt.Errorf("acmes: read key file failed, %v", keyErr)
		return
	}
	serverPEM, serverKeyPEM, serverErr := afssl.GenerateCertificate(afssl.CertificateConfig{}, afssl.WithParent(caPEM, keyPEM))
	if serverErr != nil {
		err = fmt.Errorf("acmes: generate server cert failed, %v", serverErr)
		return
	}
	clients := x509.NewCertPool()
	if !clients.AppendCertsFromPEM(caPEM) {
		err = fmt.Errorf("acmes: generate server cert failed for append client ca failed")
		return
	}
	certificate, certificateErr := tls.X509KeyPair(serverPEM, serverKeyPEM)
	if certificateErr != nil {
		err = fmt.Errorf("acmes: generate server cert failed, %v", certificateErr)
		return
	}
	config = &tls.Config{
		ClientCAs:    clients,
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	return
}
