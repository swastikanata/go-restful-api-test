package service

import (
	"context"
	"errors"
	"github.com/aronipurwanto/go-restful-api/exception"
	"github.com/aronipurwanto/go-restful-api/helper"
	"github.com/aronipurwanto/go-restful-api/model/domain"
	"github.com/aronipurwanto/go-restful-api/model/web"
	"github.com/aronipurwanto/go-restful-api/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type EmployeeServiceImpl struct {
	EmployeeRepository repository.EmployeeRepository
	Validate           *validator.Validate
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository, validate *validator.Validate) EmployeeService {
	return &EmployeeServiceImpl{
		EmployeeRepository: employeeRepository,
		Validate:           validate,
	}
}

// Create Employee
func (service *EmployeeServiceImpl) Create(ctx context.Context, request web.EmployeeCreateRequest) (web.EmployeeResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.EmployeeResponse{}, err
	}

	employee := domain.Employee{Name: request.Name}
	savedEmployee, err := service.EmployeeRepository.Save(ctx, employee)
	if err != nil {
		return web.EmployeeResponse{}, err
	}

	return helper.ToEmployeeResponse(savedEmployee), nil
}

// Update Employee
func (service *EmployeeServiceImpl) Update(ctx context.Context, request web.EmployeeUpdateRequest) (web.EmployeeResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.EmployeeResponse{}, err
	}

	employee, err := service.EmployeeRepository.FindById(ctx, request.Id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return web.EmployeeResponse{}, exception.NewNotFoundError("Employee not found")
	} else if err != nil {
		return web.EmployeeResponse{}, err
	}

	employee.Name = request.Name
	updatedEmployee, err := service.EmployeeRepository.Update(ctx, employee)
	if err != nil {
		return web.EmployeeResponse{}, err
	}

	return helper.ToEmployeeResponse(updatedEmployee), nil
}

// Delete Employee
func (service *EmployeeServiceImpl) Delete(ctx context.Context, employeeId uint64) error {
	employee, err := service.EmployeeRepository.FindById(ctx, employeeId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewNotFoundError("Employee not found")
	} else if err != nil {
		return err
	}

	return service.EmployeeRepository.Delete(ctx, employee)
}

// Find Employee By ID
func (service *EmployeeServiceImpl) FindById(ctx context.Context, employeeId uint64) (web.EmployeeResponse, error) {
	employee, err := service.EmployeeRepository.FindById(ctx, employeeId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return web.EmployeeResponse{}, exception.NewNotFoundError("Employee not found")
	} else if err != nil {
		return web.EmployeeResponse{}, err
	}

	return helper.ToEmployeeResponse(employee), nil
}

// Find All Employees
func (service *EmployeeServiceImpl) FindAll(ctx context.Context) ([]web.EmployeeResponse, error) {
	employees, err := service.EmployeeRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return helper.ToEmployeeResponses(employees), nil
}
