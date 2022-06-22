package controllers

import (
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Fieldemployee_home_redis = "LISTEMPLOYEE_BACKEND_ISBPANEL"
const Fieldemployee_frontend_redis = "LISTEMPLOYEE_FRONTEND_ISBPANEL"

func Employeehome(c *fiber.Ctx) error {
	var obj entities.Model_employee
	var arraobj []entities.Model_employee
	var objdepart entities.Model_listdepartement
	var arraobjdepart []entities.Model_listdepartement
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldemployee_home_redis)
	jsonredis := []byte(resultredis)
	listdepartement_RD, _, _, _ := jsonparser.Get(jsonredis, "listdepartement")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		employee_username, _ := jsonparser.GetString(value, "employee_username")
		employee_iddepart, _ := jsonparser.GetString(value, "employee_iddepart")
		employee_nmdepart, _ := jsonparser.GetString(value, "employee_nmdepart")
		employee_name, _ := jsonparser.GetString(value, "employee_name")
		employee_phone, _ := jsonparser.GetString(value, "employee_phone")
		employee_status, _ := jsonparser.GetString(value, "employee_status")
		employee_create, _ := jsonparser.GetString(value, "employee_create")
		employee_update, _ := jsonparser.GetString(value, "employee_update")

		obj.Employee_username = employee_username
		obj.Employee_iddepart = employee_iddepart
		obj.Employee_nmdepart = employee_nmdepart
		obj.Employee_name = employee_name
		obj.Employee_phone = employee_phone
		obj.Employee_status = employee_status
		obj.Employee_create = employee_create
		obj.Employee_update = employee_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listdepartement_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		departement_id, _ := jsonparser.GetString(value, "departement_id")
		departement_name, _ := jsonparser.GetString(value, "departement_name")

		objdepart.Departement_id = departement_id
		objdepart.Departement_name = departement_name
		arraobjdepart = append(arraobjdepart, objdepart)
	})
	if !flag {
		result, err := models.Fetch_employeeHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldemployee_home_redis, result, 60*time.Minute)
		log.Println("EMPLOYEE MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EMPLOYEE CACHE")
		return c.JSON(fiber.Map{
			"status":          fiber.StatusOK,
			"message":         "Success",
			"listdepartement": arraobjdepart,
			"record":          arraobj,
			"time":            time.Since(render_page).String(),
		})
	}
}
func EmployeeByDepart(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_employeebydepart)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	var obj entities.Model_employeebydepart
	var arraobj []entities.Model_employeebydepart
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldemployee_home_redis + "_" + client.Employee_iddepart)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		employee_username, _ := jsonparser.GetString(value, "employee_username")
		employee_name, _ := jsonparser.GetString(value, "employee_name")
		employee_deposit, _ := jsonparser.GetInt(value, "employee_deposit")
		employee_noanswer, _ := jsonparser.GetInt(value, "employee_noanswer")
		employee_reject, _ := jsonparser.GetInt(value, "employee_reject")
		employee_invalid, _ := jsonparser.GetInt(value, "employee_invalid")

		obj.Employee_username = employee_username
		obj.Employee_name = employee_name
		obj.Employee_deposit = int(employee_deposit)
		obj.Employee_noanswer = int(employee_noanswer)
		obj.Employee_reject = int(employee_reject)
		obj.Employee_invalid = int(employee_invalid)
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_employeeByDepartement(client.Employee_iddepart)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldemployee_home_redis+"_"+client.Employee_iddepart, result, 60*time.Minute)
		log.Println("EMPLOYEE BY DEPARTEMENT MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EMPLOYEE BY DEPARTEMENT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func EmployeeSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_employeesave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	// admin, username, password, iddepart, name, phone, status, sData, idrecord string
	result, err := models.Save_employee(
		client_admin,
		client.Employee_username, client.Employee_password, client.Employee_iddepart, client.Employee_name,
		client.Employee_phone, client.Employee_status, client.Sdata, client.Employee_username)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_employee()
	return c.JSON(result)
}
func _deleteredis_employee() {
	val_master := helpers.DeleteRedis(Fieldemployee_home_redis)
	log.Printf("Redis Delete BACKEND EMPLOYEE : %d", val_master)

	//CLIENT
	val_client := helpers.DeleteRedis(Fieldemployee_frontend_redis)
	log.Printf("Redis Delete FRONTEND EMPLOYEE : %d", val_client)

	val_master_departement := helpers.DeleteRedis(Fielddepartement_home_redis)
	log.Printf("Redis Delete BACKEND DEPARTEMENT : %d", val_master_departement)

	//CLIENT
	val_client_departement := helpers.DeleteRedis(Fielddepartement_frontend_redis)
	log.Printf("Redis Delete FRONTEND DEPARTEMENT : %d", val_client_departement)
}
