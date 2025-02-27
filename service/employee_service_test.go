package service

import (
	"context"
	"errors"
	"github.com/aronipurwanto/go-restful-api/model/domain"
	"github.com/aronipurwanto/go-restful-api/model/web"
	"github.com/aronipurwanto/go-restful-api/repository/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockValidator := validator.New()
	employeeService := NewEmployeeService(mockRepo, mockValidator)

	tests := []struct {
		name      string
		input     web.EmployeeCreateRequest
		mock      func()
		expect    web.EmployeeResponse
		expectErr bool
	}{
		{
			name:  "success",
			input: web.EmployeeCreateRequest{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Employee{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
			},
			expect:    web.EmployeeResponse{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			expectErr: false,
		},
		{
			name:      "validation error",
			input:     web.EmployeeCreateRequest{Name: ""},
			mock:      func() {},
			expect:    web.EmployeeResponse{},
			expectErr: true,
		},
		{
			name:  "repository error",
			input: web.EmployeeCreateRequest{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Employee{}, errors.New("database error"))
			},
			expect:    web.EmployeeResponse{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := employeeService.Create(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, resp)
			}
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEmployeeRepository(ctrl)
	employeeService := NewEmployeeService(mockRepo, validator.New())

	tests := []struct {
		name       string
		employeeId uint64
		mock       func()
		expectErr  bool
	}{
		{
			name:       "success",
			employeeId: 1,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Employee{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name:       "not found",
			employeeId: 99,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(99)).Return(domain.Employee{}, errors.New("not found"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := employeeService.Delete(context.Background(), tt.employeeId)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateEmployee(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockEmployeeRepo *mocks.MockEmployeeRepository)
		input   web.EmployeeUpdateRequest
		expects error
	}{
		{
			name: "Success",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Employee{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
				mockEmployeeRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Employee{Name: "Updated Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
			},
			input:   web.EmployeeUpdateRequest{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			expects: nil,
		},
		{
			name: "Employee Not Found",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Employee{}, errors.New("not found"))
			},
			input:   web.EmployeeUpdateRequest{Id: 1, Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			expects: errors.New("not found"),
		},
		{
			name: "Validation Error - Empty Name",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				// Tidak perlu mock FindById karena validasi gagal sebelum ke repository
			},
			input:   web.EmployeeUpdateRequest{Id: 1, Name: "", Email: "", Phone: "", DateHired: ""},
			expects: errors.New("EmployeeUpdateRequest.Name"),
		},
		{
			name: "Database Error on Update",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Employee{Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
				mockEmployeeRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Employee{}, errors.New("database error"))
			},
			input:   web.EmployeeUpdateRequest{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			expects: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEmployeeRepo := mocks.NewMockEmployeeRepository(ctrl)
			tt.mock(mockEmployeeRepo)

			service := NewEmployeeService(mockEmployeeRepo, validator.New())
			_, err := service.Update(context.Background(), tt.input)

			if tt.expects != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expects.Error()) // Alternatif untuk assert.ErrorContains
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindAllEmployees(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockEmployeeRepo *mocks.MockEmployeeRepository)
		expects []web.EmployeeResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindAll(gomock.Any()).Return([]domain.Employee{{EmployeeID: 1, Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}}, nil)
			},
			expects: []web.EmployeeResponse{{Id: 1, Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}},
			err:     nil,
		},
		{
			name: "Database Error",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expects: nil,
			err:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEmployeeRepo := mocks.NewMockEmployeeRepository(ctrl)
			tt.mock(mockEmployeeRepo)

			service := NewEmployeeService(mockEmployeeRepo, validator.New())
			result, err := service.FindAll(context.Background())
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestFindByIdEmployee(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockEmployeeRepo *mocks.MockEmployeeRepository)
		input   uint64
		expects web.EmployeeResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Employee{EmployeeID: 1, Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"}, nil)
			},
			input:   1,
			expects: web.EmployeeResponse{Id: 1, Name: "Test", Email: "test@test.com", Phone: "123456", DateHired: "01/01/2025"},
			err:     nil,
		},
		{
			name: "Not Found",
			mock: func(mockEmployeeRepo *mocks.MockEmployeeRepository) {
				mockEmployeeRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Employee{}, errors.New("not found"))
			},
			input:   1,
			expects: web.EmployeeResponse{},
			err:     errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockEmployeeRepo := mocks.NewMockEmployeeRepository(ctrl)
			tt.mock(mockEmployeeRepo)

			service := NewEmployeeService(mockEmployeeRepo, validator.New())
			result, err := service.FindById(context.Background(), tt.input)
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}
