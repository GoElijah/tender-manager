package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) RollbackBid(ctx context.Context, request gen.RollbackBidRequestObject) (gen.RollbackBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.BidId, request.Version, request.Body.Username)
	if err != nil {
		return gen.RollbackBid400JSONResponse{
			Error: "path params is invalid",
		}, nil
	}

	Bid := entity.Bid{
		ID:        request.BidId,
		Version:   request.Version,
		UpdatedBy: request.Body.Username,
	}

	userData, err := c.eapi.GetByUsername(ctx, Bid.UpdatedBy)
	if err != nil {
		return gen.RollbackBid404JSONResponse{
			Error: "User not found",
		}, nil
	}

	BidData, err := c.GetBid(ctx, Bid.ID)
	if err != nil {
		return gen.RollbackBid403JSONResponse{
			Error: "Bid not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, BidData.BidOrganization); err != nil {
		return gen.RollbackBid404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	var restoredBid entity.Bid
	restoredBid, err = c.storage.RollbackBid(ctx, Bid)
	if err != nil {
		return gen.RollbackBid500JSONResponse{
			Error: "Error rollback Bid",
		}, err
	}
	return gen.RollbackBid200JSONResponse{
		Id:          restoredBid.ID,
		Name:        restoredBid.Name,
		Description: restoredBid.Description,
		Version:     restoredBid.Version,
	}, nil
}
