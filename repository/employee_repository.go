package repository

import (
	"context"
	"github.com/aronipurwanto/go-restful-api/model/domain"
)

type EmployeeRepository interface {
	Save(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	Update(ctx context.Context, employee domain.Employee) (domain.Employee, error)
	Delete(ctx context.Context, employee domain.Employee) error
	FindById(ctx context.Context, employeeId uint64) (domain.Employee, error)
	FindAll(ctx context.Context) ([]domain.Employee, error)
}
