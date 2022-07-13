package deepl

import (
	"context"
	"net/http"
	"net/url"
)

type plan int32

const (
	Free plan = iota
	Pro
)

type Client struct {
	context.Context
	header  http.Header
	baseURL *url.URL
	plan    plan
}

func New(accessToken string, plan plan) *Client {
	ctx := context.Background()
	c := &Client{
		Context: ctx,
		header:  http.Header{},
		plan:    plan,
	}

	if plan == Free {
		c.baseURL, _ = url.Parse("https://api-free.deepl.com/v2/")
	} else if plan == Pro {
		c.baseURL, _ = url.Parse("https://api-pro.deepl.com/v2/")
	}

	c.header.Set("Authorization", "DeepL-Auth-Key " + accessToken)

	return c
}
