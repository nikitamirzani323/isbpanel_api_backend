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

const Fieldevent_home_redis = "LISTEVENT_BACKEND_ISBPANEL"

func Eventhome(c *fiber.Ctx) error {
	var obj entities.Model_event
	var arraobj []entities.Model_event
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldevent_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		event_id, _ := jsonparser.GetInt(value, "event_id")
		event_idwebagen, _ := jsonparser.GetInt(value, "event_idwebagen")
		event_name, _ := jsonparser.GetString(value, "event_name")
		event_startevent, _ := jsonparser.GetString(value, "event_startevent")
		event_endevent, _ := jsonparser.GetString(value, "event_endevent")
		event_create, _ := jsonparser.GetString(value, "event_create")
		event_update, _ := jsonparser.GetString(value, "event_update")

		obj.Event_id = int(event_id)
		obj.Event_idwebagen = int(event_idwebagen)
		obj.Event_name = event_name
		obj.Event_startevent = event_startevent
		obj.Event_endevent = event_endevent
		obj.Event_create = event_create
		obj.Event_update = event_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_event()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldevent_home_redis, result, 60*time.Minute)
		log.Println("EVENT  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EVENT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func EventSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventsave)
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

	// admin, nmevent, startevent, endevent, sData string,
	// idwebagen, idrecord int
	result, err := models.Save_event(
		client_admin,
		client.Event_name, client.Event_startevent, client.Event_endevent,
		client.Sdata, client.Event_idwebagen, client.Event_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_event()
	return c.JSON(result)
}
func _deleteredis_event() {
	val_master := helpers.DeleteRedis(Fieldevent_home_redis)
	log.Printf("Redis Delete BACKEND EVENT : %d", val_master)

}
