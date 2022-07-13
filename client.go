package deepl

import (
	"context"
	"net/http"
)

type Client struct {
	context.Context
	header http.Header
}

func New(accessToken string) *Client {
	ctx := context.Background()
	c := &Client{
		Context: ctx,
		header:  http.Header{},
	}

	c.header.Set("Authorization", "DeepL-Auth-Key "+accessToken)

	return c
}
