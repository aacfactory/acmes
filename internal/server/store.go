package server

import (
	"fmt"
	"github.com/aacfactory/acmes/internal/store"
	"net/url"
)

func createStore(p string) (v store.Store, err error) {
	u, urlErr := url.Parse(p)
	if urlErr != nil {
		err = fmt.Errorf("acmes: parse store url failed, %v", urlErr)
		return
	}
	switch u.Scheme {
	case "file":
		v, err = store.NewFileStore(p[8:])
		break
	case "oss":
		err = fmt.Errorf("acmes: store schema is not support")
		return
	default:
		err = fmt.Errorf("acmes: store schema is not support")
		return
	}
	return
}
