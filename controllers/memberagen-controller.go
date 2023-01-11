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

const Fieldmember_home_redis = "LISTMEMBER_BACKEND_ISBPANEL"
const Fieldmemberagen_home_redis = "LISTMEMBERAGEN_BACKEND_ISBPANEL"

func Memberhome(c *fiber.Ctx) error {
	var obj entities.Model_member
	var arraobj []entities.Model_member
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmember_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		member_phone, _ := jsonparser.GetString(value, "member_phone")
		member_name, _ := jsonparser.GetString(value, "member_name")
		member_create, _ := jsonparser.GetString(value, "member_create")
		member_update, _ := jsonparser.GetString(value, "member_update")

		obj.Member_phone = member_phone
		obj.Member_name = member_name
		obj.Member_create = member_create
		obj.Member_update = member_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_member()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmember_home_redis, result, 60*time.Minute)
		log.Println("MEMBER  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("MEMBER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Memberagenhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_memberagen)
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

	var obj entities.Model_memberagen
	var arraobj []entities.Model_memberagen
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmemberagen_home_redis + "_" + client.Memberagen_phone)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		memberagen_id, _ := jsonparser.GetInt(value, "memberagen_id")
		memberagen_idwebagen, _ := jsonparser.GetInt(value, "memberagen_idwebagen")
		memberagen_username, _ := jsonparser.GetString(value, "memberagen_username")
		memberagen_create, _ := jsonparser.GetString(value, "memberagen_create")
		memberagen_update, _ := jsonparser.GetString(value, "memberagen_update")

		obj.Memberagen_id = int(memberagen_id)
		obj.Memberagen_idwebagen = int(memberagen_idwebagen)
		obj.Memberagen_username = memberagen_username
		obj.Memberagen_create = memberagen_create
		obj.Memberagen_update = memberagen_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_memberagen(client.Memberagen_phone)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmemberagen_home_redis+"_"+client.Memberagen_phone, result, 60*time.Minute)
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
func MemberSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_membersave)
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

	// admin, phone, nama, sData string,
	// idrecord int
	result, err := models.Save_member(
		client_admin,
		client.Member_phone, client.Member_name, client.Sdata, client.Member_phone)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_memberagen(client.Member_phone)
	return c.JSON(result)
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
		client.Memberagen_username, client.Memberagen_phone, client.Sdata, client.Memberagen_idwebagen, client.Memberagen_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_memberagen(client.Memberagen_phone)
	return c.JSON(result)
}
func _deleteredis_memberagen(phone string) {
	val_member := helpers.DeleteRedis(Fieldmember_home_redis)
	log.Printf("Redis Delete BACKEND MEMBER  : %d", val_member)
	val_memberagen := helpers.DeleteRedis(Fieldmemberagen_home_redis + "_" + phone)
	log.Printf("Redis Delete BACKEND MEMBER AGEN : %d", val_memberagen)

}
