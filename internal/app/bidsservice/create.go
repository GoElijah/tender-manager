package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) CreateBid(ctx context.Context, request gen.CreateBidRequestObject) (gen.CreateBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.Name, request.Body.Description, request.Body.CreatorUsername, request.Body.TenderId, request.Body.TenderOrganization, request.Body.BidOrganization)
	if err != nil {
		c.logger.Error("err Validation", "err:", err)
		return gen.CreateBid400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	bid := entity.Bid{
		Name:               request.Body.Name,
		Description:        request.Body.Description,
		CreatedBy:          request.Body.CreatorUsername,
		TenderID:           request.Body.TenderId,
		TenderOrganization: request.Body.TenderOrganization,
		BidOrganization:    request.Body.BidOrganization,
	}

	bidId, err := c.storage.CreateBid(ctx, bid)
	if err != nil {
		c.logger.Error("func CreateBid", "err:", err)
		return gen.CreateBid500JSONResponse{
			Error: "Error create bid",
		}, err
	}

	return gen.CreateBid201JSONResponse{
		Id:          bidId,
		Name:        bid.Name,
		Description: bid.Description,
	}, err
}
