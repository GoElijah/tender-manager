package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) PatchTender(ctx context.Context, request gen.PatchTenderRequestObject) (gen.PatchTenderResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.EditorUsername, request.TenderId)
	if err != nil {
		return gen.PatchTender400JSONResponse{
			Error: "Editor username or tenderId is empty",
		}, nil
	}

	userData, err := c.eapi.GetByUsername(ctx, request.Body.EditorUsername)
	if err != nil {
		return gen.PatchTender404JSONResponse{
			Error: "User not found",
		}, nil
	}

	tenderData, err := c.GetTender(ctx, request.TenderId)
	if err != nil {
		return gen.PatchTender403JSONResponse{
			Error: "Tender not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, tenderData.Organization); err != nil {
		return gen.PatchTender404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	tender := entity.Tender{
		ID:          request.TenderId,
		Name:        *request.Body.Name,
		Description: *request.Body.Description,
		UpdatedBy:   request.Body.EditorUsername,
	}

	var patchedTender entity.Tender
	patchedTender, err = c.storage.PatchTender(ctx, tender)
	if err != nil {
		return gen.PatchTender500JSONResponse{
			Error: "Error patch tender",
		}, err
	}

	return gen.PatchTender200JSONResponse{
		Name:        patchedTender.Name,
		Description: patchedTender.Description,
		TenderId:    patchedTender.ID,
	}, nil
}
