package tenderservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *Client) CreateTender(ctx context.Context, request gen.CreateTenderRequestObject) (gen.CreateTenderResponseObject, error) {
	err := app.ValidateNotEmpty(request.Body.Name, request.Body.Description, request.Body.CreatorUsername, request.Body.ServiceType, request.Body.OrganizationId)
	if err != nil {
		return gen.CreateTender400JSONResponse{
			Error: "body param is invalid",
		}, nil
	}

	tender := entity.Tender{
		Name:         request.Body.Name,
		Description:  request.Body.Description,
		CreatedBy:    request.Body.CreatorUsername,
		ServiceType:  request.Body.ServiceType,
		Organization: request.Body.OrganizationId,
	}

	var tenderId string
	tenderId, err = c.storage.CreateTender(ctx, tender)
	if err != nil {
		return gen.CreateTender500JSONResponse{
			Error: "Error create tender",
		}, nil
	}

	return gen.CreateTender201JSONResponse{
		TenderId:    tenderId,
		Name:        tender.Name,
		Description: tender.Description,
	}, nil
}
