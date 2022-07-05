package controllers

import (
	"log"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Field_tafsirmimpihome_redis = "LISTTAFSIRMIMPI_BACKEND_ISBPANEL"

func Tafsirmimpihome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_tafsirmimpi)
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
	if client.Tafsirmimpi_search != "" {
		val_tafsirmimpi := helpers.DeleteRedis(Field_tafsirmimpihome_redis + "_" + client.Tafsirmimpi_search)
		log.Printf("Redis Delete BACKEND TAFSIRMIMPI : %d", val_tafsirmimpi)
	}
	var obj entities.Model_tafsirmimpi
	var arraobj []entities.Model_tafsirmimpi
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Field_tafsirmimpihome_redis + "_" + strconv.Itoa(client.Tafsirmimpi_page) + "_" + client.Tafsirmimpi_search)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		tafsirmimpi_id, _ := jsonparser.GetInt(value, "tafsirmimpi_id")
		tafsirmimpi_mimpi, _ := jsonparser.GetString(value, "tafsirmimpi_mimpi")
		tafsirmimpi_artimimpi, _ := jsonparser.GetString(value, "tafsirmimpi_artimimpi")
		tafsirmimpi_angka2d, _ := jsonparser.GetString(value, "tafsirmimpi_angka2d")
		tafsirmimpi_angka3d, _ := jsonparser.GetString(value, "tafsirmimpi_angka3d")
		tafsirmimpi_angka4d, _ := jsonparser.GetString(value, "tafsirmimpi_angka4d")
		tafsirmimpi_status, _ := jsonparser.GetString(value, "tafsirmimpi_status")
		tafsirmimpi_statuscss, _ := jsonparser.GetString(value, "tafsirmimpi_statuscss")
		tafsirmimpi_create, _ := jsonparser.GetString(value, "tafsirmimpi_create")
		tafsirmimpi_update, _ := jsonparser.GetString(value, "tafsirmimpi_update")

		obj.Tafsirmimpi_id = int(tafsirmimpi_id)
		obj.Tafsirmimpi_mimpi = tafsirmimpi_mimpi
		obj.Tafsirmimpi_artimimpi = tafsirmimpi_artimimpi
		obj.Tafsirmimpi_angka2d = tafsirmimpi_angka2d
		obj.Tafsirmimpi_angka3d = tafsirmimpi_angka3d
		obj.Tafsirmimpi_angka4d = tafsirmimpi_angka4d
		obj.Tafsirmimpi_status = tafsirmimpi_status
		obj.Tafsirmimpi_statuscss = tafsirmimpi_statuscss
		obj.Tafsirmimpi_create = tafsirmimpi_create
		obj.Tafsirmimpi_update = tafsirmimpi_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_tafsirmimpiHome(client.Tafsirmimpi_search, client.Tafsirmimpi_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Field_tafsirmimpihome_redis+"_"+strconv.Itoa(client.Tafsirmimpi_page)+"_"+client.Tafsirmimpi_search, result, 30*time.Minute)
		log.Println("TAFSIR MIMPI MYSQL")
		return c.JSON(result)
	} else {
		log.Println("TAFSIR MIMPI CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Tafsirmimpisave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_tafsirmimpisave)
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

	result, err := models.Save_tafsirmimpi(
		client_admin,
		client.Tafsirmimpi_mimpi, client.Tafsirmimpi_artimimpi, client.Tafsirmimpi_angka2d,
		client.Tafsirmimpi_angka3d, client.Tafsirmimpi_angka4d, client.Tafsirmimpi_status, client.Sdata, client.Tafsirmimpi_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_tafsirmimpi(client.Tafsirmimpi_page, client.Tafsirmimpi_search)
	val_tafsirmimpi := helpers.DeleteRedis(Field_tafsirmimpihome_redis + "_")
	log.Printf("Redis Delete BACKEND TAFSIRMIMPI : %d", val_tafsirmimpi)
	return c.JSON(result)
}
func _deleteredis_tafsirmimpi(page int, search string) {
	val_tafsirmimpi := helpers.DeleteRedis(Field_tafsirmimpihome_redis + "_" + strconv.Itoa(page) + "_" + search)
	log.Printf("Redis Delete BACKEND TAFSIRMIMPI : %d", val_tafsirmimpi)

}
