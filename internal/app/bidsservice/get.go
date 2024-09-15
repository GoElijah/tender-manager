package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) GetBid(ctx context.Context, BidID string) (entity.Bid, error) {
	var bidData entity.Bid
	bidData, err := c.storage.GetBid(ctx, BidID)
	c.logger.Error("func GetBid", "err:", err)
	if err != nil {
		return bidData, err
	}
	return bidData, nil

}

func (c *BidsClient) GetBidStatus(ctx context.Context, request gen.GetBidStatusRequestObject) (gen.GetBidStatusResponseObject, error) {
	err := app.ValidateNotEmpty(request.BidId)
	if err != nil {
		return gen.GetBidStatus400JSONResponse{
			Error: "BidId is not set in params",
		}, nil
	}
	BidData, err := c.GetBid(ctx, request.BidId)
	if err != nil {
		c.logger.Error("func GetBid", "err:", err)
		return gen.GetBidStatus500JSONResponse{
			Error: "Error getting bid status",
		}, nil
	}
	return gen.GetBidStatus200JSONResponse{
		Id:     request.BidId,
		Status: BidData.Status,
	}, nil
}

func (c *BidsClient) GetBidVersion(ctx context.Context, request gen.GetBidVersionRequestObject) (gen.GetBidVersionResponseObject, error) {
	err := app.ValidateNotEmpty(request.BidId)
	if err != nil {
		return gen.GetBidVersion400JSONResponse{
			Error: "BidId is not set in params",
		}, nil
	}
	bidData, err := c.GetBid(ctx, request.BidId)
	if err != nil {
		c.logger.Error("func GetBid", "err:", err)
		return gen.GetBidVersion500JSONResponse{
			Error: "Error closing Bid",
		}, nil
	}
	return gen.GetBidVersion200JSONResponse{
		Id:      request.BidId,
		Version: bidData.Version,
	}, nil
}
