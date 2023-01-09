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

const Fieldmemberagen_home_redis = "LISTMEMBERAGEN_BACKEND_ISBPANEL"
const Fieldmemberagen_frontend_redis = "LISTMEMBERAGEN_FRONTEND_ISBPANEL"

func Memberagenhome(c *fiber.Ctx) error {
	var obj entities.Model_memberagen
	var arraobj []entities.Model_memberagen
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmemberagen_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		memberagen_id, _ := jsonparser.GetInt(value, "memberagen_id")
		memberagen_idwebagen, _ := jsonparser.GetInt(value, "memberagen_idwebagen")
		memberagen_username, _ := jsonparser.GetString(value, "memberagen_username")
		memberagen_name, _ := jsonparser.GetString(value, "memberagen_name")
		memberagen_create, _ := jsonparser.GetString(value, "memberagen_create")
		memberagen_update, _ := jsonparser.GetString(value, "memberagen_update")

		obj.Memberagen_id = int(memberagen_id)
		obj.Memberagen_idwebagen = int(memberagen_idwebagen)
		obj.Memberagen_username = memberagen_username
		obj.Memberagen_name = memberagen_name
		obj.Memberagen_create = memberagen_create
		obj.Memberagen_update = memberagen_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_domainHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmemberagen_frontend_redis, result, 60*time.Minute)
		log.Println("MEMBER AGEN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("DOMAIN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func MemberagenSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_memberagensave)
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

	result, err := models.Save_memberagen(
		client_admin,
		client.Memberagen_username, client.Memberagen_name, client.Sdata, client.Memberagen_idwebagen, client.Memberagen_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_memberagen()
	return c.JSON(result)
}
func _deleteredis_memberagen() {
	val_master := helpers.DeleteRedis(Fieldmemberagen_home_redis)
	log.Printf("Redis Delete BACKEND MEMBER AGEN : %d", val_master)

}
