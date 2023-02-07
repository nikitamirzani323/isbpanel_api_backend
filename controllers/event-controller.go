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

const Fieldevent_home_redis = "LISTEVENT_BACKEND_ISBPANEL"
const Fieldeventdetail_home_redis = "LISTEVENTDETAIL_BACKEND_ISBPANEL"
const Fieldeventwinner_home_redis = "LISTEVENTWINNER_BACKEND_ISBPANEL"
const Fieldeventdetailgroup_home_redis = "LISTEVENTDETAILGROUP_BACKEND_ISBPANEL"

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
		event_nmwebagen, _ := jsonparser.GetString(value, "event_nmwebagen")
		event_name, _ := jsonparser.GetString(value, "event_name")
		event_startevent, _ := jsonparser.GetString(value, "event_startevent")
		event_endevent, _ := jsonparser.GetString(value, "event_endevent")
		event_mindeposit, _ := jsonparser.GetInt(value, "event_mindeposit")
		event_money_in, _ := jsonparser.GetInt(value, "event_money_in")
		event_money_out, _ := jsonparser.GetInt(value, "event_money_out")
		event_create, _ := jsonparser.GetString(value, "event_create")
		event_update, _ := jsonparser.GetString(value, "event_update")

		obj.Event_id = int(event_id)
		obj.Event_idwebagen = int(event_idwebagen)
		obj.Event_nmwebagen = event_nmwebagen
		obj.Event_name = event_name
		obj.Event_startevent = event_startevent
		obj.Event_endevent = event_endevent
		obj.Event_mindeposit = int(event_mindeposit)
		obj.Event_money_in = int(event_money_in)
		obj.Event_money_out = int(event_money_out)
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
func Eventdetailhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventdetail)
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
	var obj entities.Model_eventdetail
	var arraobj []entities.Model_eventdetail
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldeventdetail_home_redis + "_" + strconv.Itoa(client.Event_id) + "_" + strconv.Itoa(client.Event_idmemberagen))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		eventdetail_id, _ := jsonparser.GetInt(value, "eventdetail_id")
		eventdetail_deposit, _ := jsonparser.GetInt(value, "eventdetail_deposit")
		eventdetail_status, _ := jsonparser.GetString(value, "eventdetail_status")
		eventdetail_voucher, _ := jsonparser.GetString(value, "eventdetail_voucher")
		eventdetail_phone, _ := jsonparser.GetString(value, "eventdetail_phone")
		eventdetail_username, _ := jsonparser.GetString(value, "eventdetail_username")
		eventdetail_create, _ := jsonparser.GetString(value, "eventdetail_create")
		eventdetail_update, _ := jsonparser.GetString(value, "eventdetail_update")

		obj.Eventdetail_iddetail = int(eventdetail_id)
		obj.Eventdetail_deposit = int(eventdetail_deposit)
		obj.Eventdetail_voucher = eventdetail_voucher
		obj.Eventdetail_status = eventdetail_status
		obj.Eventdetail_phone = eventdetail_phone
		obj.Eventdetail_username = eventdetail_username
		obj.Eventdetail_create = eventdetail_create
		obj.Eventdetail_update = eventdetail_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetchdetail_event(client.Event_id, client.Event_idmemberagen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldeventdetail_home_redis+"_"+strconv.Itoa(client.Event_id)+"_"+strconv.Itoa(client.Event_idmemberagen), result, 60*time.Minute)
		log.Println("EVENT DETAIL  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EVENT DETAIL CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Eventdetailwinner(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventdetailwinner)
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
	var obj entities.Model_eventdetail
	var arraobj []entities.Model_eventdetail
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldeventwinner_home_redis + "_" + strconv.Itoa(client.Event_id))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		eventdetail_id, _ := jsonparser.GetInt(value, "eventdetail_id")
		eventdetail_deposit, _ := jsonparser.GetInt(value, "eventdetail_deposit")
		eventdetail_status, _ := jsonparser.GetString(value, "eventdetail_status")
		eventdetail_voucher, _ := jsonparser.GetString(value, "eventdetail_voucher")
		eventdetail_phone, _ := jsonparser.GetString(value, "eventdetail_phone")
		eventdetail_username, _ := jsonparser.GetString(value, "eventdetail_username")
		eventdetail_create, _ := jsonparser.GetString(value, "eventdetail_create")
		eventdetail_update, _ := jsonparser.GetString(value, "eventdetail_update")

		obj.Eventdetail_iddetail = int(eventdetail_id)
		obj.Eventdetail_deposit = int(eventdetail_deposit)
		obj.Eventdetail_voucher = eventdetail_voucher
		obj.Eventdetail_status = eventdetail_status
		obj.Eventdetail_phone = eventdetail_phone
		obj.Eventdetail_username = eventdetail_username
		obj.Eventdetail_create = eventdetail_create
		obj.Eventdetail_update = eventdetail_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetchdetailwinner_event(client.Event_id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldeventdetail_home_redis+"_"+strconv.Itoa(client.Event_id), result, 60*time.Minute)
		log.Println("EVENT DETAIL WINNER  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EVENT DETAIL WINNER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Eventgroupdetailhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventdetail)
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
	var obj entities.Model_eventdetailgroup
	var arraobj []entities.Model_eventdetailgroup
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldeventdetailgroup_home_redis + "_" + strconv.Itoa(client.Event_id))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		eventdetailgroup_idmember, _ := jsonparser.GetInt(value, "eventdetailgroup_idmember")
		eventdetailgroup_deposit, _ := jsonparser.GetInt(value, "eventdetailgroup_deposit")
		eventdetailgroup_voucher, _ := jsonparser.GetInt(value, "eventdetailgroup_voucher")
		eventdetailgroup_phone, _ := jsonparser.GetString(value, "eventdetailgroup_phone")
		eventdetailgroup_username, _ := jsonparser.GetString(value, "eventdetailgroup_username")

		obj.Eventdetailgroup_idmember = int(eventdetailgroup_idmember)
		obj.Eventdetailgroup_deposit = int(eventdetailgroup_deposit)
		obj.Eventdetailgroup_voucher = int(eventdetailgroup_voucher)
		obj.Eventdetailgroup_username = eventdetailgroup_username
		obj.Eventdetailgroup_phone = eventdetailgroup_phone
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetchdetailgroup_event(client.Event_id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldeventdetailgroup_home_redis+"_"+strconv.Itoa(client.Event_id), result, 60*time.Minute)
		log.Println("EVENT GROUP DETAIL  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EVENT GROUP DETAIL CACHE")
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
		client.Sdata, client.Event_idwebagen, client.Event_mindeposit, client.Event_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_event(0, 0)
	return c.JSON(result)
}
func EventDetailSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventdetailsave)
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

	result, err := models.Savedetail_event(
		client_admin,
		client.Sdata, client.Eventdetail_idevent, client.Eventdetail_idmemberagen,
		client.Eventdetail_qty, client.Eventdetail_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_event(client.Eventdetail_idevent, client.Eventdetail_idmemberagen)
	_deleteredis_event(client.Eventdetail_idevent, 0)
	return c.JSON(result)
}
func EventDetailStatusUpdate(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_eventdetailstatusupdate)
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

	result, err := models.Updatedetailstatus_event(
		client_admin, client.Eventdetail_status,
		client.Sdata, client.Eventdetail_idevent, client.Eventdetail_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_event(client.Eventdetail_idevent, 0)
	return c.JSON(result)
}
func _deleteredis_event(idevent, idmemberagen int) {
	val_master := helpers.DeleteRedis(Fieldevent_home_redis)
	log.Printf("Redis Delete BACKEND EVENT : %d", val_master)

	val_detail := helpers.DeleteRedis(Fieldeventdetail_home_redis + "_" + strconv.Itoa(idevent) + "_" + strconv.Itoa(idmemberagen))
	log.Printf("Redis Delete BACKEND EVENT DELETE : %d", val_detail)

	val_winner := helpers.DeleteRedis(Fieldeventwinner_home_redis + "_" + strconv.Itoa(idevent))
	log.Printf("Redis Delete BACKEND EVENT WINNER DELETE : %d", val_winner)

	val_groupdetail := helpers.DeleteRedis(Fieldeventdetailgroup_home_redis + "_" + strconv.Itoa(idevent))
	log.Printf("Redis Delete BACKEND EVENT GROUP DELETE : %d", val_groupdetail)
}
