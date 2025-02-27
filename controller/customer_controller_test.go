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

func setupTestAppCustomer(mockService *mocks.MockCustomerService) *fiber.App {
	app := fiber.New()
	customerController := NewCustomerController(mockService)

	api := app.Group("/api")
	customers := api.Group("/customers")
	customers.Post("/", customerController.Create)
	customers.Put("/:customerId", customerController.Update)
	customers.Delete("/:customerId", customerController.Delete)
	customers.Get("/:customerId", customerController.FindById)
	customers.Get("/", customerController.FindAll)

	return app
}

func TestCustomerController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCustomerService(ctrl)
	app := setupTestAppCustomer(mockService)

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
			name:   "Update customer - success",
			method: "PUT",
			url:    "/api/customers/1",
			body:   web.CustomerUpdateRequest{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
			setupMock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(web.CustomerResponse{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: web.WebResponse{
				Code:   http.StatusOK,
				Status: "OK",
				Data:   web.CustomerResponse{Id: 1, Name: "Updated Test", Email: "test@test.com", Phone: "123456", Address: "test street", LoyaltyPts: 1},
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
				respBody.Data = web.CustomerResponse{
					Id:         uint64(dataMap["id"].(float64)),
					Name:       dataMap["name"].(string),
					Email:      dataMap["email"].(string),
					Phone:      dataMap["phone"].(string),
					Address:    dataMap["address"].(string),
					LoyaltyPts: int(dataMap["loyalty_points"].(float64)),
				}
			}

			assert.Equal(t, tt.expectedBody, respBody)
		})
	}
}
