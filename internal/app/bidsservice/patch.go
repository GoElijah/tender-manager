package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) PatchBid(ctx context.Context, request gen.PatchBidRequestObject) (gen.PatchBidResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.EditorUsername, request.Body)
	if err != nil {
		return gen.PatchBid400JSONResponse{
			Error: "Editor username or BidId is empty",
		}, nil
	}

	userData, err := c.eapi.GetByUsername(ctx, request.Body.EditorUsername)
	if err != nil {
		return gen.PatchBid404JSONResponse{
			Error: "User not found",
		}, nil
	}

	bidData, err := c.GetBid(ctx, request.BidId)
	if err != nil {
		return gen.PatchBid403JSONResponse{
			Error: "Bid not found",
		}, nil
	}
	if err = app.ValidateUserHasAccess(userData.OrganizationID, bidData.BidOrganization); err != nil {
		return gen.PatchBid404JSONResponse{
			Error: "You now allowed for this operation",
		}, nil
	}

	bid := entity.Bid{
		ID:          request.BidId,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		UpdatedBy:   request.Body.EditorUsername,
	}

	var patchedBid entity.Bid
	patchedBid, err = c.storage.PatchBid(ctx, bid)
	if err != nil {
		return gen.PatchBid500JSONResponse{
			Error: "Error patch Bid",
		}, err
	}

	return gen.PatchBid200JSONResponse{
		Name:        patchedBid.Name,
		Description: patchedBid.Description,
		BidId:       patchedBid.ID,
	}, nil
}
