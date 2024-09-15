package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) CancelBid(ctx context.Context, request gen.CancelBidRequestObject) (gen.CancelBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.CancelUsername, request.Body.BidId)
	if err != nil {
		c.logger.Error("Validation params", "err:", err)
		return gen.CancelBid400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	bid := entity.Bid{
		ID:        request.Body.BidId,
		UpdatedBy: request.Body.CancelUsername,
	}

	userData, err := c.eapi.GetByUsername(ctx, bid.UpdatedBy)
	if err != nil {
		c.logger.Error("func GetByUsername", "err:", err)
		return gen.CancelBid404JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, bid.ID)
	if err != nil {
		c.logger.Error("func GetBid", "err:", err)
		return gen.CancelBid403JSONResponse{
			Error: "Bids not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.BidOrganization); err != nil {
		c.logger.Error("access validation", "err:", err)
		return gen.CancelBid404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	err = c.storage.CancelBid(ctx, bid)
	if err != nil {
		c.logger.Error("func CancelBid", "err:", err)
		return gen.CancelBid500JSONResponse{
			Error: "Error cancel Bid",
		}, nil
	}

	return gen.CancelBid200JSONResponse{
		Status: "Canceled",
		Name:   bidData.Name,
	}, nil
}
