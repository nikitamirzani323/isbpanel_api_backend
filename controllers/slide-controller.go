package controllers

import (
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_backend/helpers"
	"bitbucket.org/isbtotogroup/isbpanel_backend/models"
	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const Fieldslider_home_redis = "LISTSLIDER_BACKEND_ISBPANEL"

func Sliderhome(c *fiber.Ctx) error {
	var obj entities.Model_slide
	var arraobj []entities.Model_slide
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldslider_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		slide_id, _ := jsonparser.GetInt(value, "slide_id")
		slide_movieid, _ := jsonparser.GetInt(value, "slide_movieid")
		slide_movietitle, _ := jsonparser.GetString(value, "slide_movietitle")
		slide_urlimage, _ := jsonparser.GetString(value, "slide_urlimage")
		slide_position, _ := jsonparser.GetInt(value, "slide_position")
		slide_create, _ := jsonparser.GetString(value, "slide_create")

		obj.Slide_id = int(slide_id)
		obj.Slide_movieid = int(slide_movieid)
		obj.Slide_movietitle = slide_movietitle
		obj.Slide_urlimage = slide_urlimage
		obj.Slide_position = int(slide_position)
		obj.Slide_create = slide_create
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_slideHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldslider_home_redis, result, 30*time.Minute)
		log.Println("SLIDER MYSQL")
		return c.JSON(result)
	} else {
		log.Println("SLIDER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Slidersave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_slidesave)
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

	result, err := models.Save_slide(
		client_admin,
		client.Sdata, client.Slide_urlimage, client.Slide_movieid, client.Slide_position)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_slider := helpers.DeleteRedis(Fieldslider_home_redis)
	log.Printf("Redis Delete BACKEND SLIDER : %d", val_slider)
	return c.JSON(result)
}
func Sliderdelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_slidedelete)
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

	result, err := models.Delete_slider(client_admin, client.Slide_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_slider := helpers.DeleteRedis(Fieldslider_home_redis)
	log.Printf("Redis Delete BACKEND SLIDER : %d", val_slider)
	return c.JSON(result)
}
