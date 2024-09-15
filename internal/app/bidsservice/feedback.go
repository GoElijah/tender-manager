package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) Feedback(ctx context.Context, request gen.FeedbackRequestObject) (gen.FeedbackResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.PublisherUsername, request.Body.BidId, request.Body.Feedback)
	if err != nil {
		c.logger.Error("err Validation", "err:", err)
		return gen.Feedback400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	bid := entity.Bid{
		ID:        request.Body.BidId,
		UpdatedBy: request.Body.PublisherUsername,
		Feedback:  request.Body.Feedback,
	}

	userData, err := c.eapi.GetByUsername(ctx, bid.UpdatedBy)
	if err != nil {
		c.logger.Error("func GetByUsername", "err:", err)
		return gen.Feedback404JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, bid.ID)
	c.logger.Error("func GetBid", "err:", err)
	if err != nil {
		return gen.Feedback404JSONResponse{
			Error: "Bid not found",
		}, nil
	}

	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.TenderOrganization); err != nil {
		c.logger.Error("func validateAccess", "err:", err)
		return gen.Feedback403JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	err = c.storage.Feedback(ctx, bid)
	if err != nil {
		c.logger.Error("func Feedback", "err:", err)
		return gen.Feedback500JSONResponse{
			Error: "Error publish bid",
		}, err
	}

	return gen.Feedback200JSONResponse{
		BidId:    bid.ID,
		Feedback: bid.Feedback,
	}, nil
}
