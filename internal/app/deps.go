package app

import (
	"context"
	gen "tender-manager/internal/generated"
)

type StatusApp interface {
	Ping(ctx context.Context, request gen.PingRequestObject) (gen.PingResponseObject, error)
}

func (c *Status) Ping(ctx context.Context, request gen.PingRequestObject) (gen.PingResponseObject, error) {
	return gen.Ping200JSONResponse{
		Ping: c.Status,
	}, nil
}

// TODO add buildInfo data
type Status struct {
	Status string
}

func New(c *Status) *Status {
	return &Status{
		Status: "ok",
	}

}
