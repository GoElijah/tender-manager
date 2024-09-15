package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) ListBids(ctx context.Context, request gen.ListBidsRequestObject) (gen.ListBidsResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.Username, request.TenderId)
	if err != nil {
		c.logger.Error("err validation", "err:", err)
		return gen.ListBids400JSONResponse{
			Error: "Username is not set in Body",
		}, nil
	}

	var listBids []gen.Bid
	listBids, err = c.storage.ListBids(ctx, entity.Bid{TenderID: request.TenderId, CreatedBy: request.Body.Username})
	if err != nil {
		c.logger.Error("err List bids", "err:", err)
		return gen.ListBids500JSONResponse{
			Error: "Error list bids",
		}, nil
	}

	return gen.ListBids200JSONResponse{
		listBids,
	}, nil
}

func (c *BidsClient) ListMyBids(ctx context.Context, request gen.ListMyBidsRequestObject) (gen.ListMyBidsResponseObject, error) {
	err := app.ValidateNotEmpty(request.Params.Username)
	if err != nil {
		return gen.ListMyBids400JSONResponse{
			Error: "Username is not set in query params",
		}, nil
	}

	var listBids []gen.Bid
	listBids, err = c.storage.ListMyBids(ctx, *request.Params.Username)
	if err != nil {
		c.logger.Error("err ListMyBids", "err:", err)
		return gen.ListMyBids500JSONResponse{
			Error: "Error list bids",
		}, nil
	}

	return gen.ListMyBids200JSONResponse{
		listBids,
	}, nil
}
