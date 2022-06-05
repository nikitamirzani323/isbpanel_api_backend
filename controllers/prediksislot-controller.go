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

const Fieldprediksislot_home_redis = "LISTPREDIKSISLOT_BACKEND_ISBPANEL"

func Prediksislothome(c *fiber.Ctx) error {
	var obj entities.Model_prediksislot
	var arraobj []entities.Model_prediksislot
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldprediksislot_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		prediksislot_id, _ := jsonparser.GetInt(value, "prediksislot_id")
		prediksislot_nmprovider, _ := jsonparser.GetString(value, "prediksislot_nmprovider")
		prediksislot_name, _ := jsonparser.GetString(value, "prediksislot_name")
		prediksislot_prediksi, _ := jsonparser.GetInt(value, "prediksislot_prediksi")
		prediksislot_image, _ := jsonparser.GetString(value, "prediksislot_image")
		prediksislot_status, _ := jsonparser.GetString(value, "prediksislot_status")
		prediksislot_create, _ := jsonparser.GetString(value, "prediksislot_create")
		prediksislot_update, _ := jsonparser.GetString(value, "prediksislot_update")

		obj.Prediksislot_id = int(prediksislot_id)
		obj.Prediksislot_nmprovider = prediksislot_nmprovider
		obj.Prediksislot_name = prediksislot_name
		obj.Prediksislot_prediksi = int(prediksislot_prediksi)
		obj.Prediksislot_image = prediksislot_image
		obj.Prediksislot_status = prediksislot_status
		obj.Prediksislot_create = prediksislot_create
		obj.Prediksislot_update = prediksislot_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_prediksislotHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldprediksislot_home_redis, result, 60*time.Minute)
		log.Println("PREDIKSI SLOT MYSQL")
		return c.JSON(result)
	} else {
		log.Println("PREDIKSI SLOT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func PrediksislotSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_prediksislotsave)
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

	result, err := models.Save_prediksislot(
		client_admin,
		client.Prediksislot_name, client.Prediksislot_image, client.Prediksislot_status, client.Sdata,
		client.Providerslot_id, client.Prediksislot_prediksi, client.Prediksislot_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	val_master := helpers.DeleteRedis(Fieldprediksislot_home_redis)
	log.Printf("Redis Delete BACKEND PREDIKSI SLOT : %d", val_master)
	return c.JSON(result)
}
