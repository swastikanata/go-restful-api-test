// Code generated by MockGen. DO NOT EDIT.
// Source: controller/category_controller.go
//
// Generated by this command:
//
//	mockgen -source=controller/category_controller.go -destination=controller/mocks/category_controller_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	fiber "github.com/gofiber/fiber/v2"
	gomock "go.uber.org/mock/gomock"
)

// MockCategoryController is a mock of CategoryController interface.
type MockCategoryController struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryControllerMockRecorder
	isgomock struct{}
}

// MockCategoryControllerMockRecorder is the mock recorder for MockCategoryController.
type MockCategoryControllerMockRecorder struct {
	mock *MockCategoryController
}

// NewMockCategoryController creates a new mock instance.
func NewMockCategoryController(ctrl *gomock.Controller) *MockCategoryController {
	mock := &MockCategoryController{ctrl: ctrl}
	mock.recorder = &MockCategoryControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryController) EXPECT() *MockCategoryControllerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCategoryController) Create(c *fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCategoryControllerMockRecorder) Create(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCategoryController)(nil).Create), c)
}

// Delete mocks base method.
func (m *MockCategoryController) Delete(c *fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCategoryControllerMockRecorder) Delete(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCategoryController)(nil).Delete), c)
}

// FindAll mocks base method.
func (m *MockCategoryController) FindAll(c *fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCategoryControllerMockRecorder) FindAll(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCategoryController)(nil).FindAll), c)
}

// FindById mocks base method.
func (m *MockCategoryController) FindById(c *fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// FindById indicates an expected call of FindById.
func (mr *MockCategoryControllerMockRecorder) FindById(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockCategoryController)(nil).FindById), c)
}

// Update mocks base method.
func (m *MockCategoryController) Update(c *fiber.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCategoryControllerMockRecorder) Update(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCategoryController)(nil).Update), c)
}
