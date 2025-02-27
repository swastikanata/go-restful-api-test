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

func TestCreateCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerRepository(ctrl)
	mockValidator := validator.New()
	customerService := NewCustomerService(mockRepo, mockValidator)

	tests := []struct {
		name      string
		input     web.CustomerCreateRequest
		mock      func()
		expect    web.CustomerResponse
		expectErr bool
	}{
		{
			name:  "success",
			input: web.CustomerCreateRequest{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Customer{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
			},
			expect:    web.CustomerResponse{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			expectErr: false,
		},
		{
			name:      "validation error",
			input:     web.CustomerCreateRequest{Name: ""},
			mock:      func() {},
			expect:    web.CustomerResponse{},
			expectErr: true,
		},
		{
			name:  "repository error",
			input: web.CustomerCreateRequest{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Customer{}, errors.New("database error"))
			},
			expect:    web.CustomerResponse{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := customerService.Create(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, resp)
			}
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCustomerRepository(ctrl)
	customerService := NewCustomerService(mockRepo, validator.New())

	tests := []struct {
		name       string
		customerId uint64
		mock       func()
		expectErr  bool
	}{
		{
			name:       "success",
			customerId: 1,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Customer{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name:       "not found",
			customerId: 99,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(99)).Return(domain.Customer{}, errors.New("not found"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := customerService.Delete(context.Background(), tt.customerId)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockCustomerRepo *mocks.MockCustomerRepository)
		input   web.CustomerUpdateRequest
		expects error
	}{
		{
			name: "Success",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Customer{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
				mockCustomerRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Customer{Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
			},
			input:   web.CustomerUpdateRequest{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			expects: nil,
		},
		{
			name: "Customer Not Found",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Customer{}, errors.New("not found"))
			},
			input:   web.CustomerUpdateRequest{Id: 1, Name: "Tes", Email: "test@test.com", Phone: "123456", Address: "tes", LoyaltyPts: 1},
			expects: errors.New("not found"),
		},
		{
			name: "Validation Error - Empty Name",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				// Tidak perlu mock FindById karena validasi gagal sebelum ke repository
			},
			input:   web.CustomerUpdateRequest{Id: 1, Name: "", Email: "", Phone: "", Address: "", LoyaltyPts: 0},
			expects: errors.New("CustomerUpdateRequest.Name"),
		},
		{
			name: "Database Error on Update",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Customer{Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
				mockCustomerRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Customer{}, errors.New("database error"))
			},
			input:   web.CustomerUpdateRequest{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			expects: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCustomerRepo := mocks.NewMockCustomerRepository(ctrl)
			tt.mock(mockCustomerRepo)

			service := NewCustomerService(mockCustomerRepo, validator.New())
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

func TestFindAllCustomers(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockCustomerRepo *mocks.MockCustomerRepository)
		expects []web.CustomerResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindAll(gomock.Any()).Return([]domain.Customer{{CustomerID: 1, Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}}, nil)
			},
			expects: []web.CustomerResponse{{Id: 1, Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}},
			err:     nil,
		},
		{
			name: "Database Error",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expects: nil,
			err:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCustomerRepo := mocks.NewMockCustomerRepository(ctrl)
			tt.mock(mockCustomerRepo)

			service := NewCustomerService(mockCustomerRepo, validator.New())
			result, err := service.FindAll(context.Background())
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestFindByIdCustomer(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockCustomerRepo *mocks.MockCustomerRepository)
		input   uint64
		expects web.CustomerResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Customer{CustomerID: 1, Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
			},
			input:   1,
			expects: web.CustomerResponse{Id: 1, Name: "Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			err:     nil,
		},
		{
			name: "Not Found",
			mock: func(mockCustomerRepo *mocks.MockCustomerRepository) {
				mockCustomerRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Customer{}, errors.New("not found"))
			},
			input:   1,
			expects: web.CustomerResponse{},
			err:     errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCustomerRepo := mocks.NewMockCustomerRepository(ctrl)
			tt.mock(mockCustomerRepo)

			service := NewCustomerService(mockCustomerRepo, validator.New())
			result, err := service.FindById(context.Background(), tt.input)
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}
