package repository

import (
	"context"
	"errors"
	"github.com/aronipurwanto/go-restful-api/model/domain"
	"gorm.io/gorm"
)

type EmployeeRepositoryImpl struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{db: db}
}

// Save employee
func (repository *EmployeeRepositoryImpl) Save(ctx context.Context, employee domain.Employee) (domain.Employee, error) {
	if err := repository.db.WithContext(ctx).Create(&employee).Error; err != nil {
		return domain.Employee{}, err
	}
	return employee, nil
}

// Update employee
func (repository *EmployeeRepositoryImpl) Update(ctx context.Context, employee domain.Employee) (domain.Employee, error) {
	if err := repository.db.WithContext(ctx).Save(&employee).Error; err != nil {
		return domain.Employee{}, err
	}
	return employee, nil
}

// Delete employee
func (repository *EmployeeRepositoryImpl) Delete(ctx context.Context, employee domain.Employee) error {
	if err := repository.db.WithContext(ctx).Delete(&employee).Error; err != nil {
		return err
	}
	return nil
}

// FindById - Get employee by ID
func (repository *EmployeeRepositoryImpl) FindById(ctx context.Context, employeeId uint64) (domain.Employee, error) {
	var employee domain.Employee
	err := repository.db.WithContext(ctx).First(&employee, employeeId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return employee, errors.New("employee is not found")
	}
	return employee, err
}

// FindAll - Get all employees
func (repository *EmployeeRepositoryImpl) FindAll(ctx context.Context) ([]domain.Employee, error) {
	var employees []domain.Employee
	err := repository.db.WithContext(ctx).Find(&employees).Error
	return employees, err
}
