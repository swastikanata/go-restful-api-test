package web

type EmployeeCreateRequest struct {
	Name      string `validate:"required,min=1,max=100" json:"name"`
	Email     string `validate:"required,min=1,max=100" json:"email"`
	Phone     string `validate:"required,min=1,max=100" json:"phone"`
	DateHired string `validate:"required,min=0" json:"column:date_hired"`
}

type EmployeeUpdateRequest struct {
	Id        uint64 `validate:"required"`
	Name      string `validate:"required,max=200,min=1" json:"name"`
	Email     string `validate:"required,min=1,max=100" json:"email"`
	Phone     string `validate:"required,min=1,max=100" json:"phone"`
	DateHired string `validate:"required,min=0" json:"column:date_hired"`
}

type EmployeeResponse struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	DateHired string `json:"date_hired"`
}
