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

const Fieldbanner_home_redis = "LISTBANNER_BACKEND_ISBPANEL"
const Fieldbanner_frontend_redis = "LISTBANNER_FRONTEND_ISBPANEL"

func Bannerhome(c *fiber.Ctx) error {
	var obj entities.Model_banner
	var arraobj []entities.Model_banner
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldbanner_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		banner_id, _ := jsonparser.GetInt(value, "banner_id")
		banner_name, _ := jsonparser.GetString(value, "banner_name")
		banner_url, _ := jsonparser.GetString(value, "banner_url")
		banner_urlwebsite, _ := jsonparser.GetString(value, "banner_urlwebsite")
		banner_posisi, _ := jsonparser.GetString(value, "banner_posisi")
		banner_device, _ := jsonparser.GetString(value, "banner_device")
		banner_display, _ := jsonparser.GetInt(value, "banner_display")
		banner_status, _ := jsonparser.GetString(value, "banner_status")
		banner_create, _ := jsonparser.GetString(value, "banner_create")
		banner_update, _ := jsonparser.GetString(value, "banner_update")

		obj.Banner_id = int(banner_id)
		obj.Banner_name = banner_name
		obj.Banner_url = banner_url
		obj.Banner_urlwebsite = banner_urlwebsite
		obj.Banner_posisi = banner_posisi
		obj.Banner_device = banner_device
		obj.Banner_display = int(banner_display)
		obj.Banner_status = banner_status
		obj.Banner_create = banner_create
		obj.Banner_update = banner_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_bannerHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldbanner_home_redis, result, 180*time.Minute)
		log.Println("BANNER MYSQL")
		return c.JSON(result)
	} else {
		log.Println("BANNER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Bannersave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_bannerduplicate)
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

	result, err := models.Duplicate_banner(client_admin, client.Banner_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_banner()
	return c.JSON(result)
}
func _deleteredis_banner() {
	val_master := helpers.DeleteRedis(Fieldbanner_home_redis)
	log.Printf("Redis Delete BACKEND BANNER : %d", val_master)

	//CLIENT
	val_client_banner := helpers.DeleteRedis(Fieldbanner_frontend_redis)
	log.Printf("Redis Delete CLIENT BANNER : %d", val_client_banner)
}
