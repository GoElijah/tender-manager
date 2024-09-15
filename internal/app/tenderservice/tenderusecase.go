package tenderservice

import (
	"context"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

type tenderStorage interface {
	GetTender(ctx context.Context, tenderID string) (entity.Tender, error)
	CreateTender(ctx context.Context, t entity.Tender) (string, error)
	PublishTender(ctx context.Context, t entity.Tender) (string, error)
	CloseTender(ctx context.Context, t entity.Tender) (string, error)
	PatchTender(ctx context.Context, t entity.Tender) (entity.Tender, error)
	ListTenders(ctx context.Context, t entity.Tender) ([]gen.Tender, error)
	ListMyTenders(ctx context.Context, t entity.Tender) ([]gen.Tender, error)
	RollbackTender(ctx context.Context, t entity.Tender) (entity.Tender, error)
}
