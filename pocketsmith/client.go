package pocketsmith

import "errors"

var ErrNotFound = errors.New("not found")

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}
