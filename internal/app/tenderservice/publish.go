package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) PublishTender(ctx context.Context, request gen.PublishTenderRequestObject) (gen.PublishTenderResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.PublisherUsername, request.Body.TenderId)
	if err != nil {
		return gen.PublishTender400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	tender := entity.Tender{
		ID:        request.Body.TenderId,
		UpdatedBy: request.Body.PublisherUsername,
	}

	userData, err := c.eapi.GetByUsername(ctx, tender.UpdatedBy)
	if err != nil {
		return gen.PublishTender404JSONResponse{
			Error: "User not found",
		}, nil
	}

	tenderData, err := c.GetTender(ctx, tender.ID)
	if err != nil {
		return gen.PublishTender404JSONResponse{
			Error: "Tender not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, tenderData.Organization); err != nil {
		return gen.PublishTender404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	tenderId, err := c.storage.PublishTender(ctx, tender)
	if err != nil {
		return gen.PublishTender500JSONResponse{
			Error: "Error publish tender",
		}, nil
	}
	return gen.PublishTender200JSONResponse{
		TenderId: tenderId,
		Status:   "PUBLISHED",
	}, nil
}
