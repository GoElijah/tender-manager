package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) CloseTender(ctx context.Context, request gen.CloseTenderRequestObject) (gen.CloseTenderResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.CloserUsername, request.Body.TenderId)
	if err != nil {
		return gen.CloseTender400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	tender := entity.Tender{
		ID:        request.Body.TenderId,
		UpdatedBy: request.Body.CloserUsername,
	}

	userData, err := c.eapi.GetByUsername(ctx, tender.UpdatedBy)
	if err != nil {
		return gen.CloseTender404JSONResponse{
			Error: "User not found",
		}, nil
	}

	tenderData, err := c.GetTender(ctx, tender.ID)
	if err != nil {
		return gen.CloseTender403JSONResponse{
			Error: "Tender not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, tenderData.Organization); err != nil {
		return gen.CloseTender404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	tenderId, err := c.storage.CloseTender(ctx, tender)
	if err != nil {
		return gen.CloseTender500JSONResponse{
			Error: "Error close tender",
		}, err
	}
	return gen.CloseTender200JSONResponse{
		TenderId: tenderId,
		Status:   "CLOSED",
	}, nil
}
