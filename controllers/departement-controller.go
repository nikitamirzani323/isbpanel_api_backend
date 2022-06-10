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

const Fielddepartement_home_redis = "LISTDEPARTEMENT_BACKEND_ISBPANEL"
const Fielddepartement_frontend_redis = "LISTDEPARTEMENT_FRONTEND_ISBPANEL"

func Departementhome(c *fiber.Ctx) error {
	var obj entities.Model_departement
	var arraobj []entities.Model_departement
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fielddepartement_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		departement_id, _ := jsonparser.GetString(value, "departement_id")
		departement_name, _ := jsonparser.GetString(value, "departement_name")
		departement_create, _ := jsonparser.GetString(value, "departement_create")
		departement_update, _ := jsonparser.GetString(value, "departement_update")

		obj.Departement_id = departement_id
		obj.Departement_name = departement_name
		obj.Departement_create = departement_create
		obj.Departement_update = departement_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_departementHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fielddepartement_home_redis, result, 60*time.Minute)
		log.Println("DEPARTEMENT MYSQL")
		return c.JSON(result)
	} else {
		log.Println("DEPARTEMENT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func DepartementSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_departementsave)
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

	// admin, nmdepartement, status, sData string, idrecord int
	result, err := models.Save_departement(
		client_admin,
		client.Departement_name, client.Sdata, client.Departement_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_departement()
	return c.JSON(result)
}
func _deleteredis_departement() {
	val_master := helpers.DeleteRedis(Fielddepartement_home_redis)
	log.Printf("Redis Delete BACKEND DEPARTEMENT : %d", val_master)

	//CLIENT
	val_client := helpers.DeleteRedis(Fielddepartement_frontend_redis)
	log.Printf("Redis Delete FRONTEND DEPARTEMENT : %d", val_client)

	val_master_employee := helpers.DeleteRedis(Fieldemployee_home_redis)
	log.Printf("Redis Delete BACKEND EMPLOYEE : %d", val_master_employee)

	//CLIENT
	val_client_employee := helpers.DeleteRedis(Fieldemployee_frontend_redis)
	log.Printf("Redis Delete FRONTEND EMPLOYEE : %d", val_client_employee)

}
