package server

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/aacfactory/acmes/internal/store"
	"github.com/aacfactory/logs"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"golang.org/x/sync/singleflight"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestParam struct {
	Domain string `json:"domain"`
}

type Handler struct {
	log     logs.Logger
	email   string
	acme    *lego.Client
	stores  store.Store
	barrier *singleflight.Group
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler.log.DebugEnabled() {
		handler.log.Debug().Message(fmt.Sprintf("%s %s", request.Method, request.URL.String()))
	}
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}
	contentType := request.Header.Get("Content-Type")
	if contentType != "application/acme" {
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}
	param := &RequestParam{}
	body, bodyErr := ioutil.ReadAll(request.Body)
	if bodyErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	paramErr := json.Unmarshal(body, param)
	if paramErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	requestPath := request.URL.Path
	var cert *store.Certificate
	var err error
	switch requestPath {
	case "/obtain":
		cert, err = handler.obtain(handler.email, param.Domain)
	case "/renew":
		cert, err = handler.renew(handler.email, param.Domain)
	default:
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(fmt.Sprintf("{\"cause\": \"%s\"}", err.Error())))
		return
	}
	result, certErr := json.Marshal(cert)
	if certErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(fmt.Sprintf("{\"cause\": \"%s\"}", certErr.Error())))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	_, _ = writer.Write(result)
}

func (handler *Handler) obtain(email string, domain string) (v *store.Certificate, err error) {
	if handler.log.DebugEnabled() {
		handler.log.Debug().Message(fmt.Sprintf("begin obtain %s", domain))
	}
	key := fmt.Sprintf("obtain:%s:%s", email, domain)
	result, doErr, _ := handler.barrier.Do(key, func() (v interface{}, handleErr error) {
		cert, hasCert, getErr := handler.stores.GetUserCertificate(context.TODO(), email, domain)
		if getErr != nil {
			handleErr = getErr
			return
		}
		if hasCert {
			v = cert
			return
		}
		request := certificate.ObtainRequest{
			Domains: []string{domain},
			Bundle:  true,
		}
		certificates, obtainErr := handler.acme.Certificate.Obtain(request)
		if obtainErr != nil {
			handleErr = obtainErr
			return
		}
		cert, handleErr = handler.handleCertificates(certificates)
		if handleErr != nil {
			return
		}

		saveErr := handler.stores.SaveUserCertificate(context.TODO(), email, domain, cert)
		if saveErr != nil {
			handleErr = saveErr
			return
		}
		v = cert
		return
	})
	handler.barrier.Forget(key)
	if doErr != nil {
		if handler.log.DebugEnabled() {
			handler.log.Debug().Cause(doErr).Message(fmt.Sprintf("obtain %s failed", domain))
		}
		err = fmt.Errorf("acmes: obtain failed, %v", doErr)
		return
	}
	if handler.log.DebugEnabled() {
		handler.log.Debug().Message(fmt.Sprintf("obtain %s succeed", domain))
	}
	v = result.(*store.Certificate)
	return
}

func (handler *Handler) renew(email string, domain string) (v *store.Certificate, err error) {
	if handler.log.DebugEnabled() {
		handler.log.Debug().Message(fmt.Sprintf("begin renew %s", domain))
	}
	key := fmt.Sprintf("renew:%s:%s", email, domain)
	result, doErr, _ := handler.barrier.Do(key, func() (v interface{}, handleErr error) {
		cert, hasCert, getErr := handler.stores.GetUserCertificate(context.TODO(), email, domain)
		if getErr != nil {
			handleErr = getErr
			return
		}
		if hasCert {
			handleErr = fmt.Errorf("not obtained")
			return
		}
		if cert.NotAfter.After(time.Now()) {
			v = cert
			return
		}
		resource := &certificate.Resource{}
		resourceErr := json.Unmarshal(cert.Resource, resource)
		if resourceErr != nil {
			handleErr = resourceErr
			return
		}
		certificates, renewErr := handler.acme.Certificate.Renew(*resource, true, true, "")
		if renewErr != nil {
			handleErr = renewErr
			return
		}
		cert, handleErr = handler.handleCertificates(certificates)
		if handleErr != nil {
			return
		}
		saveErr := handler.stores.SaveUserCertificate(context.TODO(), email, domain, cert)
		if saveErr != nil {
			handleErr = saveErr
			return
		}
		v = cert
		return
	})
	handler.barrier.Forget(key)
	if doErr != nil {
		if handler.log.DebugEnabled() {
			handler.log.Debug().Cause(doErr).Message(fmt.Sprintf("renew %s failed", domain))
		}
		err = fmt.Errorf("acmes: renew failed, %v", doErr)
		return
	}
	if handler.log.DebugEnabled() {
		handler.log.Debug().Cause(doErr).Message(fmt.Sprintf("renew %s succeed", domain))
	}
	v = result.(*store.Certificate)
	return
}

func (handler *Handler) handleCertificates(certificates *certificate.Resource) (v *store.Certificate, err error) {
	resp, getErr := http.Get(certificates.CertStableURL)
	if getErr != nil {
		err = fmt.Errorf("get cert from %s failed, %v", certificates.CertStableURL, getErr)
		return
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("get cert from %s failed, %v", certificates.CertStableURL, string(body))
		return
	}
	certPEM, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		err = readErr
		return
	}
	certBlock, _ := pem.Decode(certPEM)
	cert0, parseCertificateErr := x509.ParseCertificate(certBlock.Bytes)
	if parseCertificateErr != nil {
		err = parseCertificateErr
		return
	}
	renewAT := cert0.NotAfter.Local()
	keyPEM := certificates.PrivateKey
	resource, resourceErr := json.Marshal(certificates)
	if resourceErr != nil {
		err = resourceErr
		return
	}
	v = &store.Certificate{
		Resource: resource,
		Cert:     certPEM,
		Key:      keyPEM,
		NotAfter: renewAT,
	}
	return
}
