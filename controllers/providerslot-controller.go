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

const Fieldproviderslot_home_redis = "LISTPROVIDERSLOT_BACKEND_ISBPANEL"

func Providerslothome(c *fiber.Ctx) error {
	var obj entities.Model_providerslot
	var arraobj []entities.Model_providerslot
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldproviderslot_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		providerslot_id, _ := jsonparser.GetInt(value, "providerslot_id")
		providerslot_name, _ := jsonparser.GetString(value, "providerslot_name")
		providerslot_display, _ := jsonparser.GetInt(value, "providerslot_display")
		providerslot_counter, _ := jsonparser.GetInt(value, "providerslot_counter")
		providerslot_totalgameslot, _ := jsonparser.GetInt(value, "providerslot_totalgameslot")
		providerslot_image, _ := jsonparser.GetString(value, "providerslot_image")
		providerslot_slug, _ := jsonparser.GetString(value, "providerslot_slug")
		providerslot_title, _ := jsonparser.GetString(value, "providerslot_title")
		providerslot_descp, _ := jsonparser.GetString(value, "providerslot_descp")
		providerslot_status, _ := jsonparser.GetString(value, "providerslot_status")
		providerslot_create, _ := jsonparser.GetString(value, "providerslot_create")
		providerslot_update, _ := jsonparser.GetString(value, "providerslot_update")

		obj.Providerslot_id = int(providerslot_id)
		obj.Providerslot_name = providerslot_name
		obj.Providerslot_display = int(providerslot_display)
		obj.Providerslot_counter = int(providerslot_counter)
		obj.Providerslot_totalgameslot = int(providerslot_totalgameslot)
		obj.Providerslot_image = providerslot_image
		obj.Providerslot_slug = providerslot_slug
		obj.Providerslot_title = providerslot_title
		obj.Providerslot_descp = providerslot_descp
		obj.Providerslot_status = providerslot_status
		obj.Providerslot_create = providerslot_create
		obj.Providerslot_update = providerslot_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_providerslotHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldproviderslot_home_redis, result, 60*time.Minute)
		log.Println("PROVIDER SLOT MYSQL")
		return c.JSON(result)
	} else {
		log.Println("PROVIDER SLOT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func ProviderslotSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_providerslotsave)
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

	result, err := models.Save_providerslot(
		client_admin,
		client.Providerslot_name, client.Providerslot_image, client.Providerslot_slug, client.Providerslot_title,
		client.Providerslot_descp, client.Providerslot_status, client.Sdata, client.Providerslot_display, client.Providerslot_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_providerslot(client.Providerslot_slug)
	return c.JSON(result)
}
func _deleteredis_providerslot(slug string) {
	val_master := helpers.DeleteRedis(Fieldproviderslot_home_redis)
	log.Printf("Redis Delete BACKEND PROVIDER SLOT : %d", val_master)

	val_client_providerslot := helpers.DeleteRedis("LISTPROVIDERSLOT_FRONTEND_ISBPANEL")
	log.Printf("Redis Delete client PREDIKSI SLOT : %d", val_client_providerslot)

	val_client_providerslot_slug := helpers.DeleteRedis("LISTPROVIDERSLOT_FRONTEND_ISBPANEL_" + slug)
	log.Printf("Redis Delete client PREDIKSI SLOT SLUG : %d", val_client_providerslot_slug)

	val_client_prediksislot := helpers.DeleteRedis("LISTPREDIKSISLOT_FRONTEND_ISBPANEL")
	log.Printf("Redis Delete client PREDIKSI SLOT : %d", val_client_prediksislot)

	val_client_prediksislot_slug := helpers.DeleteRedis("LISTPREDIKSISLOT_FRONTEND_ISBPANEL_" + slug)
	log.Printf("Redis Delete client PREDIKSI SLOT SLUG : %d", val_client_prediksislot_slug)
}
