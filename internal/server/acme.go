package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/aacfactory/acmes/internal/store"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
)

func createAcme(email string, dnsProvider string, stores store.Store) (client *lego.Client, err error) {
	user, hasUser, getUserErr := stores.GetUser(context.TODO(), email)
	if getUserErr != nil {
		err = fmt.Errorf("acmes: create acme client failed, %v", getUserErr)
		return
	}
	if !hasUser {
		key, keyErr := rsa.GenerateKey(rand.Reader, 2048)
		if keyErr != nil {
			err = fmt.Errorf("acmes: create acme client failed, create user private key failed, %v", keyErr)
			return
		}
		user = &store.User{
			Email:    email,
			Resource: nil,
			Key: pem.EncodeToMemory(&pem.Block{
				Type:    "RSA PRIVATE KEY",
				Headers: nil,
				Bytes:   x509.MarshalPKCS1PrivateKey(key),
			}),
		}
	}
	config := lego.NewConfig(user)
	client, err = lego.NewClient(config)
	if err != nil {
		err = fmt.Errorf("acmes: create acme failed, %v", err)
		return
	}
	provider, providerErr := dns.NewDNSChallengeProviderByName(dnsProvider)
	if providerErr != nil {
		err = fmt.Errorf("acmes: create acme failed, %v", providerErr)
		return
	}
	setProviderErr := client.Challenge.SetDNS01Provider(provider)
	if setProviderErr != nil {
		err = fmt.Errorf("acmes: create acme failed, %v", setProviderErr)
		return
	}
	if !hasUser {
		userRegistration, registerErr := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if registerErr != nil {
			err = fmt.Errorf("acmes: create acme failed, %v", registerErr)
			return
		}
		userRegistrationContent, userRegistrationErr := json.Marshal(userRegistration)
		if userRegistrationErr != nil {
			err = fmt.Errorf("acmes: create acme failed, %v", registerErr)
			return
		}
		user.Resource = userRegistrationContent
		saveErr := stores.SaveUser(context.TODO(), user)
		if saveErr != nil {
			err = fmt.Errorf("acmes: create acme failed, %v", saveErr)
		}
		return
	}
	return
}
