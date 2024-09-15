package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) RejectBid(ctx context.Context, request gen.RejectBidRequestObject) (gen.RejectBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.Username, request.Body.BidId)
	if err != nil {
		c.logger.Error("err validation", "err:", err)
		return gen.RejectBid400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	userData, err := c.eapi.GetByUsername(ctx, request.Body.Username)
	if err != nil {
		return gen.RejectBid400JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, request.Body.BidId)
	if err != nil {
		return gen.RejectBid400JSONResponse{
			Error: "Bid not found",
		}, nil
	}

	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.TenderOrganization); err != nil {
		c.logger.Error("err validation", "err:", err)
		return gen.RejectBid400JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	submittedBid, err := c.storage.RejectBid(ctx, entity.Bid{UpdatedBy: request.Body.Username, ID: request.Body.BidId})
	if err != nil {
		c.logger.Error("func rejectBid", "err:", err)

		return gen.RejectBid500JSONResponse{
			Error: "Error submit bid",
		}, err
	}

	return gen.RejectBid200JSONResponse{
		Status: submittedBid.Status,
		Name:   submittedBid.Name,
		Votes:  submittedBid.Votes,
	}, nil
}
