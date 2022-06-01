package store

import (
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/go-acme/lego/v4/registration"
	"golang.org/x/net/context"
	"time"
)

type Store interface {
	GetUser(ctx context.Context, email string) (user *User, has bool, err error)
	SaveUser(ctx context.Context, user *User) (err error)
	GetUserCertificate(ctx context.Context, email string, domain string) (cert *Certificate, has bool, err error)
	SaveUserCertificate(ctx context.Context, email string, domain string, cert *Certificate) (err error)
}

type User struct {
	Email    string `json:"email"`
	Resource []byte `json:"resource"`
	Key      []byte `json:"key"`
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	r := &registration.Resource{}
	rErr := json.Unmarshal(u.Resource, r)
	if rErr != nil {
		return nil
	}
	return r
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	block, _ := pem.Decode(u.Key)
	key, parseKeyErr := x509.ParsePKCS1PrivateKey(block.Bytes)
	if parseKeyErr != nil {
		return nil
	}
	return key
}

type Certificate struct {
	Resource []byte    `json:"resource"`
	Cert     []byte    `json:"cert"`
	Key      []byte    `json:"key"`
	NotAfter time.Time `json:"notAfter"`
}
