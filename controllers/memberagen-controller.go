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

		var objwebsiteagen entities.Model_memberagen
		var arraobjwebsiteagen []entities.Model_memberagen
		record_memberagen_RD, _, _, _ := jsonparser.Get(value, "member_agen")
		jsonparser.ArrayEach(record_memberagen_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			memberagen_idwebagen, _ := jsonparser.GetInt(value, "memberagen_idwebagen")
			memberagen_username, _ := jsonparser.GetString(value, "memberagen_username")
			memberagen_website, _ := jsonparser.GetString(value, "memberagen_website")

			objwebsiteagen.Memberagen_idwebagen = int(memberagen_idwebagen)
			objwebsiteagen.Memberagen_username = memberagen_username
			objwebsiteagen.Memberagen_website = memberagen_website
			arraobjwebsiteagen = append(arraobjwebsiteagen, objwebsiteagen)
		})

		obj.Member_phone = member_phone
		obj.Member_name = member_name
		obj.Member_agen = arraobjwebsiteagen
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
		client.Member_phone, client.Member_name, string(client.Member_listagen), client.Sdata, client.Member_phone)
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

func _deleteredis_memberagen(phone string) {
	val_member := helpers.DeleteRedis(Fieldmember_home_redis)
	log.Printf("Redis Delete BACKEND MEMBER  : %d", val_member)

}
