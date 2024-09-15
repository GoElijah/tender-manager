package bidsservice

import (
	"context"
	"tender-manager/internal/app"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
)

func (c *BidsClient) ListFeedback(ctx context.Context, request gen.ListFeedbackRequestObject) (gen.ListFeedbackResponseObject, error) {
	err := app.ValidateNotEmpty(request.TenderId, request.Params.AuthorUsername, request.Params.OrganizationId)
	if err != nil {
		c.logger.Error("err validation", "err:", err)
		return gen.ListFeedback400JSONResponse{
			Error: "Error in params setting",
		}, nil
	}

	var bidInfo = entity.Bid{
		TenderOrganization: *request.Params.OrganizationId,
		TenderID:           request.TenderId,
		CreatedBy:          *request.Params.AuthorUsername,
	}
	var ListFeedback []gen.Feedback
	ListFeedback, err = c.storage.ListFeedback(ctx, bidInfo)
	if err != nil {
		c.logger.Error("err List feedbacks", "err:", err)
		return gen.ListFeedback500JSONResponse{
			Error: "Error list feedbacks",
		}, nil
	}

	return gen.ListFeedback200JSONResponse{
		ListFeedback,
	}, nil
}
