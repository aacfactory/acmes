package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/aacfactory/afssl"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func New(caPEM []byte, caKeyPem []byte, host string) (v *Client, err error) {
	host = strings.TrimSpace(host)
	if host == "" {
		err = fmt.Errorf("acmes: host is empty")
		return
	}
	config := afssl.CertificateConfig{
		Country:            "",
		Province:           "",
		City:               "",
		Organization:       "",
		OrganizationalUnit: "",
		CommonName:         "acmes",
		IPs:                nil,
		Emails:             nil,
		DNSNames:           nil,
	}
	cert, key, genSslErr := afssl.GenerateCertificate(config, afssl.WithExpirationDays(365), afssl.WithParent(caPEM, caKeyPem))
	if genSslErr != nil {
		err = fmt.Errorf("acmes: generate client tls failed, %v", genSslErr)
		return
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		err = fmt.Errorf("acmes: generate client cert failed for append root ca failed")
		return
	}
	certificate, certificateErr := tls.X509KeyPair(cert, key)
	if certificateErr != nil {
		err = fmt.Errorf("acmes: generate client cert failed, %v", certificateErr)
		return
	}
	tlsConfig := &tls.Config{
		RootCAs:            roots,
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: true,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	v = &Client{
		host:       host,
		httpClient: httpClient,
	}
	return
}

type Client struct {
	host       string
	httpClient *http.Client
}

func (c *Client) Obtain(ctx context.Context, domain string) (config *tls.Config, cancelAutoRenew func(), err error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		err = fmt.Errorf("acmes: obtain failed for domain is empty")
		return
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	u := url.URL{}
	u.Scheme = "https"
	u.Host = c.host
	u.Path = "/obtain"
	resp, postErr := c.httpClient.Post(u.String(), "application/acme", bytes.NewReader([]byte(fmt.Sprintf("{\"domain\":\"%s\"}", domain))))
	if postErr != nil {
		err = fmt.Errorf("acmes: obtain failed, %v", postErr)
		return
	}
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		err = fmt.Errorf("acmes: obtain failed, %v", bodyErr)
		return
	}
	if resp.StatusCode != 200 {
		handleErr := &HandleError{}
		decodeErr := json.Unmarshal(body, handleErr)
		if decodeErr != nil {
			err = fmt.Errorf("acmes: obtain failed, %v", decodeErr)
			return
		}
		err = fmt.Errorf("acmes: obtain failed, %v", handleErr.Cause)
		return
	}
	cert := &Certificate{}
	decodeErr := json.Unmarshal(body, cert)
	if decodeErr != nil {
		err = fmt.Errorf("acmes: obtain failed, %v", decodeErr)
		return
	}
	certificate, certificateErr := tls.X509KeyPair(cert.Cert, cert.Key)
	if certificateErr != nil {
		err = fmt.Errorf("acmes: obtain failed, %v", certificateErr)
		return
	}
	config = &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	cancelAutoRenew, err = c.autoRenew(ctx, domain, config, cert.NotAfter)
	return
}

func (c *Client) autoRenew(ctx context.Context, domain string, config *tls.Config, notAfter time.Time) (cancelAutoRenew func(), err error) {
	ctx, cancelAutoRenew = context.WithCancel(ctx)
	go func(ctx context.Context, domain string, config *tls.Config, c *Client, notAfter time.Time) {
		for {
			closed := false
			select {
			case <-ctx.Done():
				closed = true
				break
			case <-time.After(notAfter.Sub(time.Now())):
				notAfter = c.renew(domain, config)
				break
			}
			if closed {
				break
			}
		}
	}(ctx, domain, config, c, notAfter)
	return
}

func (c *Client) renew(domain string, config *tls.Config) (notAfter time.Time) {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = c.host
	u.Path = "/renew"
	resp, postErr := c.httpClient.Post(u.String(), "application/acme", bytes.NewReader([]byte(fmt.Sprintf("{\"domain\":\"%s\"}", domain))))
	if postErr != nil {
		notAfter = time.Now().Add(60 * time.Second)
		return
	}
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		notAfter = time.Now().Add(60 * time.Second)
		return
	}
	if resp.StatusCode != 200 {
		handleErr := &HandleError{}
		decodeErr := json.Unmarshal(body, handleErr)
		if decodeErr != nil {
			notAfter = time.Now().Add(60 * time.Second)
			return
		}
		notAfter = time.Now().Add(60 * time.Second)
		return
	}
	cert := &Certificate{}
	decodeErr := json.Unmarshal(body, cert)
	if decodeErr != nil {
		notAfter = time.Now().Add(60 * time.Second)
		return
	}
	certificate, certificateErr := tls.X509KeyPair(cert.Cert, cert.Key)
	if certificateErr != nil {
		notAfter = time.Now().Add(60 * time.Second)
		return
	}
	config.Certificates[0] = certificate
	notAfter = cert.NotAfter
	return
}
