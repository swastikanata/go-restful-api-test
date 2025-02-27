package controller

import (
	"errors"
	"github.com/aronipurwanto/go-restful-api/controller/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCategoryController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCtrl := mocks.NewMockCategoryController(ctrl)
	app := fiber.New()

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Create category - success",
			method: "POST",
			url:    "/category",
			body:   `{"name": "Electronics"}`,
			setupMock: func() {
				mockCtrl.EXPECT().Create(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:   "Create category - failure",
			method: "POST",
			url:    "/category",
			body:   `{"name": ""}`,
			setupMock: func() {
				mockCtrl.EXPECT().Create(gomock.Any()).Return(errors.New("invalid input"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:   "Update category - success",
			method: "PUT",
			url:    "/category/1",
			body:   `{"name": "Updated"}`,
			setupMock: func() {
				mockCtrl.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Find all categories",
			method: "GET",
			url:    "/categories",
			setupMock: func() {
				mockCtrl.EXPECT().FindAll(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Find category by ID - success",
			method: "GET",
			url:    "/category/1",
			setupMock: func() {
				mockCtrl.EXPECT().FindById(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Delete category - success",
			method: "DELETE",
			url:    "/category/1",
			setupMock: func() {
				mockCtrl.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
