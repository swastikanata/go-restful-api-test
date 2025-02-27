package controller

import (
	"bytes"
	"encoding/json"
	"github.com/aronipurwanto/go-restful-api/model/web"
	"github.com/aronipurwanto/go-restful-api/service/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestAppProduct(mockService *mocks.MockProductService) *fiber.App {
	app := fiber.New()
	productController := NewProductController(mockService)

	api := app.Group("/api")
	products := api.Group("/products")
	products.Post("/", productController.Create)
	products.Put("/:productId", productController.Update)
	products.Delete("/:productId", productController.Delete)
	products.Get("/:productId", productController.FindById)
	products.Get("/", productController.FindAll)

	return app
}

func TestProductController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockProductService(ctrl)
	app := setupTestAppProduct(mockService)

	tests := []struct {
		name           string
		method         string
		url            string
		body           interface{}
		setupMock      func()
		expectedStatus int
		expectedBody   web.WebResponse
	}{
		{
			name:   "Update product - success",
			method: "PUT",
			url:    "/api/products/1",
			body:   web.ProductUpdateRequest{Id: 1, Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			setupMock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(web.ProductResponse{Id: 1, Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: web.WebResponse{
				Code:   http.StatusOK,
				Status: "OK",
				Data:   web.ProductResponse{Id: 1, Name: "Updated Test", Description: "Test", Price: 1, StockQty: 1, CategoryID: 1, SKU: "test", TaxRate: 1.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody web.WebResponse
			json.NewDecoder(resp.Body).Decode(&respBody)

			if dataMap, ok := respBody.Data.(map[string]interface{}); ok {
				respBody.Data = web.ProductResponse{
					Id:          uint64(dataMap["id"].(float64)),
					Name:        dataMap["name"].(string),
					Description: dataMap["description"].(string),
					Price:       dataMap["price"].(float64),
					StockQty:    int(dataMap["stock_qty"].(float64)),
					CategoryID:  uint64(dataMap["category_id"].(float64)),
					SKU:         dataMap["sku"].(string),
					TaxRate:     dataMap["tax_rate"].(float64),
				}
			}

			assert.Equal(t, tt.expectedBody, respBody)
		})
	}
}
