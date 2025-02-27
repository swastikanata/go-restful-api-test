package main

import (
	"github.com/aronipurwanto/go-restful-api/app"
	"github.com/aronipurwanto/go-restful-api/controller"
	"github.com/aronipurwanto/go-restful-api/helper"
	"github.com/aronipurwanto/go-restful-api/model/domain"
	"github.com/aronipurwanto/go-restful-api/repository"
	"github.com/aronipurwanto/go-restful-api/service"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	server := fiber.New()

	// Initialize Database
	db := app.NewDB()

	// Run Auto Migration (Opsional, bisa dihapus jika tidak diperlukan)
	err := db.AutoMigrate(&domain.Category{}, &domain.Customer{}, &domain.Product{}, &domain.Employee{})
	helper.PanicIfError(err)

	// Initialize Validator
	validate := validator.New()

	// Initialize Repository, Service, and Controller
	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository, validate)
	categoryController := controller.NewCategoryController(categoryService)

	employeeRepository := repository.NewCategoryRepository(db)
	employeeService := service.NewCategoryService(employeeRepository, validate)
	employeeController := controller.NewCategoryController(employeeService)

	productRepository := repository.NewCategoryRepository(db)
	productService := service.NewCategoryService(productRepository, validate)
	productController := controller.NewCategoryController(productService)

	customerRepository := repository.NewCategoryRepository(db)
	customerService := service.NewCategoryService(customerRepository, validate)
	customerController := controller.NewCategoryController(customerService)

	// Setup Routes
	app.NewRouter(server, categoryController, customerController, productController, employeeController)

	// Start Server
	log.Println("Server running on port 8080")
	err = server.Listen(":8080")
	helper.PanicIfError(err)
}
