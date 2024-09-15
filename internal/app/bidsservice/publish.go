package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) PublishBid(ctx context.Context, request gen.PublishBidRequestObject) (gen.PublishBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.PublisherUsername, request.Body.BidId)
	if err != nil {
		c.logger.Error("err Validation", "err:", err)
		return gen.PublishBid400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	bid := entity.Bid{
		ID:        request.Body.BidId,
		UpdatedBy: request.Body.PublisherUsername,
	}

	userData, err := c.eapi.GetByUsername(ctx, bid.UpdatedBy)
	if err != nil {
		c.logger.Error("func GetByUsername", "err:", err)
		return gen.PublishBid404JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, bid.ID)
	c.logger.Error("func GetBid", "err:", err)
	if err != nil {
		return gen.PublishBid404JSONResponse{
			Error: "Bid not found",
		}, nil
	}

	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.BidOrganization); err != nil {
		c.logger.Error("func validateAccess", "err:", err)
		return gen.PublishBid403JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	err = c.storage.PublishBid(ctx, bid)
	if err != nil {
		c.logger.Error("func PublishBid", "err:", err)
		return gen.PublishBid500JSONResponse{
			Error: "Error publish bid",
		}, err
	}

	return gen.PublishBid200JSONResponse{
		Status: "PUBLISHED",
		Id:     bid.ID,
	}, nil
}
