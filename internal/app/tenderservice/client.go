package tenderservice

import eapi "tender-manager/internal/app/employeeservice"

type Client struct {
	storage tenderStorage
	eapi    eapi.Client
}

func New(t tenderStorage, e eapi.Client) *Client {
	return &Client{
		storage: t,
		eapi:    e,
	}
}
