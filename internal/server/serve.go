package server

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/sync/singleflight"
	slog "log"
	"net/http"
	"strings"
)

type options struct {
	port         int
	ca           string
	key          string
	level        string
	logFormatter string
	store        string
	email        string
	provider     string
}

func serve(opt options) (err error) {
	port := opt.port
	if port < 1 {
		port = 443
	}
	log, logErr := createLog(opt.level, opt.logFormatter)
	if logErr != nil {
		err = fmt.Errorf("acmes: serve failed, %v", logErr)
		return
	}
	tlsConfig, tlsErr := createTLSConfig(opt.ca, opt.key)
	if tlsErr != nil {
		err = fmt.Errorf("acmes: serve failed, %v", tlsErr)
		return
	}
	stores, storeErr := createStore(opt.store)
	if storeErr != nil {
		err = fmt.Errorf("acmes: serve failed, %v", storeErr)
		return
	}

	client, clientErr := createAcme(opt.email, opt.provider, stores)
	if clientErr != nil {
		err = fmt.Errorf("acmes: serve failed, %v", clientErr)
		return
	}

	ln, lnErr := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if lnErr != nil {
		err = fmt.Errorf("acmes: serve failed, %v", lnErr)
		return
	}
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: &Handler{
			log:     log,
			email:   strings.TrimSpace(opt.email),
			acme:    client,
			stores:  stores,
			barrier: &singleflight.Group{},
		},
		TLSConfig: tlsConfig,
		ErrorLog: slog.New(
			&writer{
				core: log,
			},
			"",
			slog.LstdFlags,
		),
	}
	if log.DebugEnabled() {
		log.Debug().Message(fmt.Sprintf("serve at :%d", port))
	}
	err = srv.Serve(ln)
	if err != nil {
		err = fmt.Errorf("acmes: serve failed, %v", err)
		return
	}
	return
}
