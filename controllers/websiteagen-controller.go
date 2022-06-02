package controllers

import (
	"log"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_backend/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_backend/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Fieldwebsiteagen_home_redis = "LISTWEBSITEAGEN_BACKEND_ISBPANEL"

func Websiteagenhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_websiteagen)
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
	if client.Websiteagen_search != "" {
		val_news := helpers.DeleteRedis(Fieldwebsiteagen_home_redis + "_" + strconv.Itoa(client.Websiteagen_page) + "_" + client.Websiteagen_search)
		log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	}
	var obj entities.Model_websiteagen
	var arraobj []entities.Model_websiteagen
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldwebsiteagen_home_redis + "_" + strconv.Itoa(client.Websiteagen_page) + "_" + client.Websiteagen_search)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		websiteagen_id, _ := jsonparser.GetInt(value, "websiteagen_id")
		websiteagen_name, _ := jsonparser.GetString(value, "websiteagen_name")
		websiteagen_status, _ := jsonparser.GetString(value, "websiteagen_status")
		websiteagen_create, _ := jsonparser.GetString(value, "websiteagen_create")
		websiteagen_update, _ := jsonparser.GetString(value, "websiteagen_update")

		obj.Websiteagen_id = int(websiteagen_id)
		obj.Websiteagen_name = websiteagen_name
		obj.Websiteagen_status = websiteagen_status
		obj.Websiteagen_create = websiteagen_create
		obj.Websiteagen_update = websiteagen_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_websiteagenHome(client.Websiteagen_search, client.Websiteagen_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldwebsiteagen_home_redis+"_"+strconv.Itoa(client.Websiteagen_page)+"_"+client.Websiteagen_search, result, 1*time.Hour)
		log.Println("WEBSITEAGEN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("WEBSITEAGEN CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     message_RD,
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func Websiteagensave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_websiteagensave)
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

	result, err := models.Save_websiteagen(
		client_admin, client.Sdata,
		client.Websiteagen_name, client.Websiteagen_status, client.Websiteagen_idrecord)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_WEBSITEAGEN := helpers.DeleteRedis(Fieldwebsiteagen_home_redis + "_" + strconv.Itoa(client.Websiteagen_page) + "_" + client.Websiteagen_search)
	log.Printf("Redis Delete BACKEND WEBSITEAGEN : %d", val_WEBSITEAGEN)
	return c.JSON(result)
}
