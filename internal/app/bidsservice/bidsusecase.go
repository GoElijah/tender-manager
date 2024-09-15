package bidsservice

import (
	"context"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

type bidsStorage interface {
	CreateBid(ctx context.Context, b entity.Bid) (string, error)
	GetBid(ctx context.Context, bidId string) (entity.Bid, error)
	PublishBid(ctx context.Context, b entity.Bid) error
	CancelBid(ctx context.Context, b entity.Bid) error
	PatchBid(ctx context.Context, b entity.Bid) (entity.Bid, error)
	SubmitBid(ctx context.Context, b entity.Bid) (entity.Bid, error)
	RejectBid(ctx context.Context, b entity.Bid) (entity.Bid, error)
	ApproveBid(ctx context.Context, b entity.Bid) (entity.Bid, error)
	ListBids(ctx context.Context, b entity.Bid) ([]gen.Bid, error)
	ListMyBids(ctx context.Context, username string) ([]gen.Bid, error)
	RollbackBid(ctx context.Context, b entity.Bid) (entity.Bid, error)
	Feedback(ctx context.Context, b entity.Bid) error
	ListFeedback(ctx context.Context, b entity.Bid) ([]gen.Feedback, error)
}
