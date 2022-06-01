package client

import "time"

type HandleError struct {
	Cause string `json:"cause"`
}

type Certificate struct {
	Resource []byte    `json:"resource"`
	Cert     []byte    `json:"cert"`
	Key      []byte    `json:"key"`
	NotAfter time.Time `json:"notAfter"`
}
