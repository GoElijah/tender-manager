package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

const quorumVotes = 3

func (c *BidsClient) SubmitBid(ctx context.Context, request gen.SubmitBidRequestObject) (gen.SubmitBidResponseObject, error) {

	bidReqInfo := entity.Bid{ID: request.Body.BidId, UpdatedBy: request.Body.Username}

	err := app.ValidateNotEmpty(request.Body.Username, request.Body.BidId)
	if err != nil {
		c.logger.Error("err vailidate", "err:", err)
		return gen.SubmitBid400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	userData, err := c.eapi.GetByUsername(ctx, request.Body.Username)
	if err != nil {
		return gen.SubmitBid400JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, request.Body.BidId)
	if err != nil {
		return gen.SubmitBid400JSONResponse{
			Error: "Bid not found",
		}, nil
	}

	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.TenderOrganization); err != nil {
		c.logger.Error("access validation", "err:", err)
		return gen.SubmitBid400JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	submittedBid, err := c.storage.SubmitBid(ctx, bidReqInfo)
	if err != nil {
		return gen.SubmitBid500JSONResponse{
			Error: "Error submit bid",
		}, nil
	}

	if submittedBid.Votes >= quorumVotes {
		submittedBid, err = c.storage.ApproveBid(ctx, bidReqInfo)
		if err != nil {
			c.logger.Error("err ApproveBid", "err:", err)
			return gen.SubmitBid500JSONResponse{
				Error: "Error approve bid",
			}, nil
		}
	}

	return gen.SubmitBid200JSONResponse{
		Status: submittedBid.Status,
		Name:   submittedBid.Name,
		Votes:  submittedBid.Votes,
	}, nil
}
