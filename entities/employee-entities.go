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
	Employee_deposit  int    `json:"employee_deposit"`
	Employee_noanswer int    `json:"employee_noanswer"`
	Employee_reject   int    `json:"employee_reject"`
	Employee_invalid  int    `json:"employee_invalid"`
}
type Model_employeebysalesperform struct {
	Sales_deposit      int         `json:"sales_deposit"`
	Sales_depositsum   float32     `json:"sales_depositsum"`
	Sales_noanswer     int         `json:"sales_noanswer"`
	Sales_reject       int         `json:"sales_reject"`
	Sales_invalid      int         `json:"sales_invalid"`
	Sales_listdeposit  interface{} `json:"sales_listdeposit"`
	Sales_listnoanswer interface{} `json:"sales_listnoanswer"`
	Sales_listinvalid  interface{} `json:"sales_listinvalid"`
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
type Controller_employeebysalesperform struct {
	Page               string `json:"page" validate:"required"`
	Employee_iddepart  string `json:"employee_iddepart" validate:"required"`
	Employee_username  string `json:"employee_username" validate:"required"`
	Employee_startdate string `json:"employee_startdate"`
	Employee_enddate   string `json:"employee_enddate"`
}
