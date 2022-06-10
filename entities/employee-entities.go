package entities

type Model_employee struct {
	Employee_username string `json:"employee_username"`
	Employee_iddepart string `json:"employee_iddepart"`
	Employee_nmdepart string `json:"employee_nmdepart"`
	Employee_name     string `json:"employee_name"`
	Employee_phone    string `json:"employee_phone"`
	Employee_status   string `json:"employee_status"`
	Employee_create   string `json:"employee_create"`
	Employee_update   string `json:"employee_update"`
}
type Model_employeebydepart struct {
	Employee_username string `json:"employee_username"`
	Employee_name     string `json:"employee_name"`
}

type Controller_employeesave struct {
	Page              string `json:"page" validate:"required"`
	Sdata             string `json:"sdata" validate:"required"`
	Employee_username string `json:"employee_username" validate:"required"`
	Employee_password string `json:"employee_password"`
	Employee_iddepart string `json:"employee_iddepart" validate:"required"`
	Employee_name     string `json:"employee_name" validate:"required"`
	Employee_phone    string `json:"employee_phone" validate:"required"`
	Employee_status   string `json:"employee_status" validate:"required"`
}
type Controller_employeebydepart struct {
	Page              string `json:"page" validate:"required"`
	Employee_iddepart string `json:"employee_iddepart" validate:"required"`
}
