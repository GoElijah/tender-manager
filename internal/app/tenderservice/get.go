package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) GetTender(ctx context.Context, tenderID string) (entity.Tender, error) {
	var tenderData entity.Tender
	tenderData, err := c.storage.GetTender(ctx, tenderID)
	if err != nil {
		return tenderData, err
	}
	return tenderData, nil

}

func (c *Client) GetTenderStatus(ctx context.Context, request gen.GetTenderStatusRequestObject) (gen.GetTenderStatusResponseObject, error) {
	err := app.ValidateNotEmpty(request.TenderId)
	if err != nil {
		return gen.GetTenderStatus400JSONResponse{
			Error: "TenderId is not set in params",
		}, nil
	}
	tenderData, err := c.GetTender(ctx, request.TenderId)
	if err != nil {
		return gen.GetTenderStatus500JSONResponse{
			Error: "Error closing tender",
		}, nil
	}
	return gen.GetTenderStatus200JSONResponse{
		Id:     request.TenderId,
		Status: tenderData.Status,
	}, nil
}

func (c *Client) GetTenderVersion(ctx context.Context, request gen.GetTenderVersionRequestObject) (gen.GetTenderVersionResponseObject, error) {
	err := app.ValidateNotEmpty(request.TenderId)
	if err != nil {
		return gen.GetTenderVersion400JSONResponse{
			Error: "TenderId is not set in params",
		}, nil
	}
	tenderData, err := c.GetTender(ctx, request.TenderId)
	if err != nil {
		return gen.GetTenderVersion500JSONResponse{
			Error: "Error closing tender",
		}, nil
	}
	return gen.GetTenderVersion200JSONResponse{
		Id:      request.TenderId,
		Version: tenderData.Version,
	}, nil
}
