package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) ListTenders(ctx context.Context, request gen.ListTendersRequestObject) (gen.ListTendersResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.Username)
	if err != nil {
		return gen.ListTenders400JSONResponse{
			Error: "Username is not set in Body",
		}, nil
	}

	var emptyStr string
	serviceType := request.Params.ServiceType
	if serviceType == nil {
		serviceType = &emptyStr
	}

	userData, err := c.eapi.GetByUsername(ctx, request.Body.Username)
	if err != nil {
		return gen.ListTenders400JSONResponse{
			Error: "not found",
		}, nil
	}

	tender := entity.Tender{
		ServiceType:  *serviceType,
		Organization: userData.OrganizationID}

	var listTenders []gen.Tender
	listTenders, err = c.storage.ListTenders(ctx, tender)
	if err != nil {
		return gen.ListTenders500JSONResponse{
			Error: "Error list tender",
		}, nil
	}

	return gen.ListTenders200JSONResponse{
		listTenders,
	}, nil
}

func (c *Client) ListMyTenders(ctx context.Context, request gen.ListMyTendersRequestObject) (gen.ListMyTendersResponseObject, error) {
	err := app.ValidateNotEmpty(*request.Params.Username)
	if err != nil {
		return gen.ListMyTenders400JSONResponse{
			Error: "Username is not set in query params",
		}, nil
	}

	tender := entity.Tender{
		CreatedBy: *request.Params.Username}

	var listTenders []gen.Tender
	listTenders, err = c.storage.ListMyTenders(ctx, tender)
	if err != nil {
		return gen.ListMyTenders500JSONResponse{
			Error: "Error list tender",
		}, nil
	}

	return gen.ListMyTenders200JSONResponse{
		listTenders,
	}, nil
}
