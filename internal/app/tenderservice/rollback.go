package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) RollbackTender(ctx context.Context, request gen.RollbackTenderRequestObject) (gen.RollbackTenderResponseObject, error) {
	err := app.ValidateNotEmpty(request.TenderId, request.Version, request.Body.Username)
	if err != nil {
		return gen.RollbackTender400JSONResponse{
			Error: "path params is invalid",
		}, nil
	}

	tender := entity.Tender{
		ID:        request.TenderId,
		Version:   request.Version,
		UpdatedBy: request.Body.Username,
	}

	userData, err := c.eapi.GetByUsername(ctx, tender.UpdatedBy)
	if err != nil {
		return gen.RollbackTender404JSONResponse{
			Error: "User not found",
		}, nil
	}

	tenderData, err := c.GetTender(ctx, tender.ID)
	if err != nil {
		return gen.RollbackTender403JSONResponse{
			Error: "Tender not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, tenderData.Organization); err != nil {
		return gen.RollbackTender404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	var restoredTender entity.Tender
	restoredTender, err = c.storage.RollbackTender(ctx, tender)
	if err != nil {
		return gen.RollbackTender500JSONResponse{
			Error: "Error rollback tender",
		}, err
	}
	return gen.RollbackTender200JSONResponse{
		Id:          restoredTender.ID,
		Name:        restoredTender.Name,
		Description: restoredTender.Description,
		Version:     restoredTender.Version,
	}, nil
}
