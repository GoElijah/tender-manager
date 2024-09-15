package employeeservice

import (
	"context"
	"tender-manager/internal/entity"
)

type employeeStorage interface {
	GetByUsername(ctx context.Context, username string) (entity.Employee, error)
}
