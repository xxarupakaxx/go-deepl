package deepl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type plan int32

const (
	Free plan = iota
	Pro
)

type authKey string

const auth authKey = "auth"

type Client struct {
	ctx     context.Context
	header  http.Header
	baseURL *url.URL
	plan    plan
}

func New(ctx context.Context, accessToken string, plan plan) *Client {
	c := &Client{
		ctx:    ctx,
		header: http.Header{},
		plan:   plan,
	}

	if plan == Free {
		c.baseURL, _ = url.Parse("https://api-free.deepl.com/v2/")
	} else if plan == Pro {
		c.baseURL, _ = url.Parse("https://api.deepl.com/v2/")
	}

	c.header.Set("Authorization", "DeepL-Auth-Key "+accessToken)
	c.ctx = context.WithValue(c.ctx, auth, accessToken)

	return c
}

func (c *Client) GetAuthKey() (string, error) {
	v, ok := c.ctx.Value(auth).(string)
	if !ok {
		return "", fmt.Errorf("failed to get authkey")
	}

	return v, nil
}
