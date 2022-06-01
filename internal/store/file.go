package store

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func NewFileStore(rootDir string) (v Store, err error) {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		err = fmt.Errorf("acmes: new file store failed for root dir is empty")
		return
	}
	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		err = fmt.Errorf("acmes: new file store failed for root dir is invalid")
		return
	}
	fs := &FileStore{
		mutex:   sync.Mutex{},
		rootDir: "",
	}
	if !fs.pathExist(rootDir) {
		mkdirErr := os.MkdirAll(rootDir, 0600)
		if mkdirErr != nil {
			err = fmt.Errorf("acmes: new file store failed for create root dir failed, %v", mkdirErr)
			return
		}
	}
	fs.rootDir = rootDir
	v = fs
	return
}

type FileStore struct {
	mutex   sync.Mutex
	rootDir string
}

func (fs *FileStore) GetUser(_ context.Context, email string) (user *User, has bool, err error) {
	fs.mutex.Lock()
	fs.mutex.Unlock()
	email = strings.TrimSpace(email)
	if email == "" {
		err = fmt.Errorf("acmes: get user failed for email is empty")
		return
	}
	userDir := filepath.Join(fs.rootDir, email)
	if !fs.pathExist(userDir) {
		return
	}
	resPath := filepath.Join(userDir, "user.json")
	if !fs.pathExist(resPath) {
		return
	}
	resource, readResourceErr := ioutil.ReadFile(resPath)
	if readResourceErr != nil {
		err = fmt.Errorf("acmes: get user failed, %v", readResourceErr)
		return
	}
	keyPath := filepath.Join(userDir, "key.pem")
	if !fs.pathExist(resPath) {
		return
	}
	key, readKeyErr := ioutil.ReadFile(keyPath)
	if readResourceErr != nil {
		err = fmt.Errorf("acmes: get user failed, %v", readKeyErr)
		return
	}
	user = &User{
		Email:    email,
		Resource: resource,
		Key:      key,
	}
	has = true
	return
}

func (fs *FileStore) SaveUser(_ context.Context, user *User) (err error) {
	fs.mutex.Lock()
	fs.mutex.Unlock()
	email := user.Email
	userDir := filepath.Join(fs.rootDir, email)
	if !fs.pathExist(userDir) {
		mkdirErr := os.MkdirAll(userDir, 0600)
		if mkdirErr != nil {
			err = fmt.Errorf("acmes: save user failed for create user dir failed, %v", mkdirErr)
			return
		}
	}
	resPath := filepath.Join(userDir, "user.json")
	saveResErr := ioutil.WriteFile(resPath, user.Resource, 0600)
	if saveResErr != nil {
		err = fmt.Errorf("acmes: save user failed, %v", saveResErr)
		return
	}
	keyPath := filepath.Join(userDir, "key.pem")
	saveKeyErr := ioutil.WriteFile(keyPath, user.Key, 0600)
	if saveKeyErr != nil {
		err = fmt.Errorf("acmes: save user failed, %v", saveKeyErr)
		return
	}
	return
}

func (fs *FileStore) GetUserCertificate(_ context.Context, email string, domain string) (cert *Certificate, has bool, err error) {
	fs.mutex.Lock()
	fs.mutex.Unlock()
	email = strings.TrimSpace(email)
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return
	}
	if strings.Contains(domain, "*") {
		domain = strings.ReplaceAll(domain, "*", "[x]")
	}
	domainDir := filepath.Join(fs.rootDir, email, domain)
	if !fs.pathExist(domainDir) {
		return
	}
	certPath := filepath.Join(domainDir, "cert.pem")
	if !fs.pathExist(certPath) {
		return
	}
	certPem, certReadErr := ioutil.ReadFile(certPath)
	if certReadErr != nil {
		err = fmt.Errorf("acmes: get user certificate failed, %v", certReadErr)
		return
	}
	keyPath := filepath.Join(domainDir, "key.pem")
	if !fs.pathExist(keyPath) {
		return
	}
	keyPem, keyReadErr := ioutil.ReadFile(keyPath)
	if keyReadErr != nil {
		err = fmt.Errorf("acmes: get user certificate failed, %v", keyReadErr)
		return
	}
	resPath := filepath.Join(domainDir, "cert.json")
	if !fs.pathExist(resPath) {
		return
	}
	res, resReadErr := ioutil.ReadFile(resPath)
	if resReadErr != nil {
		err = fmt.Errorf("acmes: get user certificate failed, %v", resReadErr)
		return
	}
	notAfterPath := filepath.Join(domainDir, "expiration.txt")
	if !fs.pathExist(notAfterPath) {
		return
	}
	notAfterContent, notAfterReadErr := ioutil.ReadFile(notAfterPath)
	if notAfterReadErr != nil {
		err = fmt.Errorf("acmes: get user certificate failed, %v", notAfterReadErr)
		return
	}
	notAfter, notAfterErr := time.Parse(time.RFC3339, strings.TrimSpace(string(notAfterContent)))
	if notAfterErr != nil {
		err = fmt.Errorf("acmes: get user certificate failed, %v", notAfterErr)
		return
	}
	cert = &Certificate{
		Resource: res,
		Cert:     certPem,
		Key:      keyPem,
		NotAfter: notAfter,
	}
	has = true
	return
}

func (fs *FileStore) SaveUserCertificate(_ context.Context, email string, domain string, cert *Certificate) (err error) {
	fs.mutex.Lock()
	fs.mutex.Unlock()
	email = strings.TrimSpace(email)
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return
	}
	if strings.Contains(domain, "*") {
		domain = strings.ReplaceAll(domain, "*", "[x]")
	}
	domainDir := filepath.Join(fs.rootDir, email, domain)
	if !fs.pathExist(domainDir) {
		mkdirErr := os.MkdirAll(domainDir, 0600)
		if mkdirErr != nil {
			err = fmt.Errorf("acmes: save user certificate failed for create user certificate dir failed, %v", mkdirErr)
			return
		}
	}
	certPath := filepath.Join(domainDir, "cert.pem")
	saveCertErr := ioutil.WriteFile(certPath, cert.Cert, 0600)
	if saveCertErr != nil {
		err = fmt.Errorf("acmes: save user certificate failed, %v", saveCertErr)
		return
	}
	keyPath := filepath.Join(domainDir, "key.pem")
	saveKeyErr := ioutil.WriteFile(keyPath, cert.Key, 0600)
	if saveKeyErr != nil {
		err = fmt.Errorf("acmes: save user certificate failed, %v", saveKeyErr)
		return
	}
	resPath := filepath.Join(domainDir, "cert.json")
	saveResErr := ioutil.WriteFile(resPath, cert.Resource, 0600)
	if saveResErr != nil {
		err = fmt.Errorf("acmes: save user certificate failed, %v", saveResErr)
		return
	}
	notAfterPath := filepath.Join(domainDir, "expiration.txt")
	block, _ := pem.Decode(cert.Cert)
	certificate, parseCertificateErr := x509.ParseCertificate(block.Bytes)
	if parseCertificateErr != nil {
		err = fmt.Errorf("acmes: save user certificate failed, %v", parseCertificateErr)
		return
	}
	notAfter := certificate.NotAfter.Format(time.RFC3339)
	saveNotAfterErr := ioutil.WriteFile(notAfterPath, []byte(notAfter), 0600)
	if saveResErr != nil {
		err = fmt.Errorf("acmes: save user certificate failed, %v", saveNotAfterErr)
		return
	}
	return
}

func (fs *FileStore) pathExist(v string) (ok bool) {
	_, err := os.Stat(v)
	if err == nil {
		ok = true
		return
	}
	ok = !os.IsNotExist(err)
	return
}
