package gomap

import (
	"github.com/cwinters8/gomap/client"
)

func NewClient(jmapSessionURL, bearerToken string) (*client.Client, error) {
	return client.NewClient(jmapSessionURL, bearerToken)
}
