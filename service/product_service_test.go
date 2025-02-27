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

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	mockValidator := validator.New()
	productService := NewProductService(mockRepo, mockValidator)

	tests := []struct {
		name      string
		input     web.ProductCreateRequest
		mock      func()
		expect    web.ProductResponse
		expectErr bool
	}{
		{
			name:  "success",
			input: web.ProductCreateRequest{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Product{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
			},
			expect:    web.ProductResponse{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			expectErr: false,
		},
		{
			name:      "validation error",
			input:     web.ProductCreateRequest{Name: ""},
			mock:      func() {},
			expect:    web.ProductResponse{},
			expectErr: true,
		},
		{
			name:  "repository error",
			input: web.ProductCreateRequest{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			mock: func() {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Product{}, errors.New("database error"))
			},
			expect:    web.ProductResponse{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := productService.Create(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, resp)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProductRepository(ctrl)
	productService := NewProductService(mockRepo, validator.New())

	tests := []struct {
		name      string
		productId uint64
		mock      func()
		expectErr bool
	}{
		{
			name:      "success",
			productId: 1,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Product{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name:      "not found",
			productId: 99,
			mock: func() {
				mockRepo.EXPECT().FindById(gomock.Any(), uint64(99)).Return(domain.Product{}, errors.New("not found"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := productService.Delete(context.Background(), tt.productId)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockProductRepo *mocks.MockProductRepository)
		input   web.ProductUpdateRequest
		expects error
	}{
		{
			name: "Success",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Product{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
				mockProductRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Product{Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
			},
			input:   web.ProductUpdateRequest{Id: 1, Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			expects: nil,
		},
		{
			name: "Product Not Found",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Product{}, errors.New("not found"))
			},
			input:   web.ProductUpdateRequest{Id: 1, Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			expects: errors.New("not found"),
		},
		{
			name: "Validation Error - Empty Name",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				// Tidak perlu mock FindById karena validasi gagal sebelum ke repository
			},
			input:   web.ProductUpdateRequest{Id: 1, Name: "", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			expects: errors.New("ProductUpdateRequest.Name"),
		},
		{
			name: "Database Error on Update",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindById(gomock.Any(), uint64(1)).
					Return(domain.Product{Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
				mockProductRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(domain.Product{}, errors.New("database error"))
			},
			input:   web.ProductUpdateRequest{Id: 1, Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			expects: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProductRepo := mocks.NewMockProductRepository(ctrl)
			tt.mock(mockProductRepo)

			service := NewProductService(mockProductRepo, validator.New())
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

func TestFindAllProducts(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockProductRepo *mocks.MockProductRepository)
		expects []web.ProductResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindAll(gomock.Any()).Return([]domain.Product{{ProductID: 1, Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}}, nil)
			},
			expects: []web.ProductResponse{{Id: 1, Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0}},
			err:     nil,
		},
		{
			name: "Database Error",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expects: nil,
			err:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProductRepo := mocks.NewMockProductRepository(ctrl)
			tt.mock(mockProductRepo)

			service := NewProductService(mockProductRepo, validator.New())
			result, err := service.FindAll(context.Background())
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestFindByIdProduct(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockProductRepo *mocks.MockProductRepository)
		input   uint64
		expects web.ProductResponse
		err     error
	}{
		{
			name: "Success",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Product{ProductID: 1, Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryId: 1, SKU: "test", TaxRate: 1.0}, nil)
			},
			input:   1,
			expects: web.ProductResponse{Id: 1, Name: "Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			err:     nil,
		},
		{
			name: "Not Found",
			mock: func(mockProductRepo *mocks.MockProductRepository) {
				mockProductRepo.EXPECT().FindById(gomock.Any(), uint64(1)).Return(domain.Product{}, errors.New("not found"))
			},
			input:   1,
			expects: web.ProductResponse{},
			err:     errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockProductRepo := mocks.NewMockProductRepository(ctrl)
			tt.mock(mockProductRepo)

			service := NewProductService(mockProductRepo, validator.New())
			result, err := service.FindById(context.Background(), tt.input)
			assert.Equal(t, tt.expects, result)
			assert.Equal(t, tt.err, err)
		})
	}
}
