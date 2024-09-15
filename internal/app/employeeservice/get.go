package employeeservice

import (
	"context"
	"tender-manager/internal/entity"
)

func (c *Client) GetByUsername(ctx context.Context, username string) (entity.Employee, error) {
	userInfo, err := c.storage.GetByUsername(ctx, username)
	if err != nil {
		return entity.Employee{}, err
	}
	return userInfo, nil
}
