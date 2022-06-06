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

const Field_home_redis = "LISTPASARAN_BACKEND_ISBPANEL"
const Field_keluaran_redis = "LISTKELUARAN_BACKEND_ISBPANEL"
const Field_prediksi_redis = "LISTPREDIKSI_BACKEND_ISBPANEL"

const Field_client_pasaran_redis = "LISTPASARAN_FRONTEND_ISBPANEL"
const Field_client_keluaran_redis = "LISTKELUARAN_FRONTEND_ISBPANEL"

func Pasaranhome(c *fiber.Ctx) error {
	var obj entities.Model_pasaran
	var arraobj []entities.Model_pasaran
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Field_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		pasaran_id, _ := jsonparser.GetString(value, "pasaran_id")
		pasaran_name, _ := jsonparser.GetString(value, "pasaran_name")
		pasaran_url, _ := jsonparser.GetString(value, "pasaran_url")
		pasaran_diundi, _ := jsonparser.GetString(value, "pasaran_diundi")
		pasaran_jamjadwal, _ := jsonparser.GetString(value, "pasaran_jamjadwal")
		pasaran_display, _ := jsonparser.GetInt(value, "pasaran_display")
		pasaran_slug, _ := jsonparser.GetString(value, "pasaran_slug")
		pasaran_meta_title, _ := jsonparser.GetString(value, "pasaran_meta_title")
		pasaran_meta_descp, _ := jsonparser.GetString(value, "pasaran_meta_descp")
		pasaran_status, _ := jsonparser.GetString(value, "pasaran_status")
		pasaran_statuscss, _ := jsonparser.GetString(value, "pasaran_statuscss")
		pasaran_keluaran, _ := jsonparser.GetString(value, "pasaran_keluaran")
		pasaran_prediksi, _ := jsonparser.GetString(value, "pasaran_prediksi")
		pasaran_create, _ := jsonparser.GetString(value, "pasaran_create")
		pasaran_update, _ := jsonparser.GetString(value, "pasaran_update")

		obj.Pasaran_id = pasaran_id
		obj.Pasaran_name = pasaran_name
		obj.Pasaran_url = pasaran_url
		obj.Pasaran_diundi = pasaran_diundi
		obj.Pasaran_jamjadwal = pasaran_jamjadwal
		obj.Pasaran_display = int(pasaran_display)
		obj.Pasaran_slug = pasaran_slug
		obj.Pasaran_meta_title = pasaran_meta_title
		obj.Pasaran_meta_descp = pasaran_meta_descp
		obj.Pasaran_status = pasaran_status
		obj.Pasaran_statuscss = pasaran_statuscss
		obj.Pasaran_keluaran = pasaran_keluaran
		obj.Pasaran_prediksi = pasaran_prediksi
		obj.Pasaran_create = pasaran_create
		obj.Pasaran_update = pasaran_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_pasaranHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Field_home_redis, result, 180*time.Minute)
		log.Println("PASARAN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("PASARAN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Pasaransave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_pasaransave)
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

	result, err := models.Save_pasaran(
		client_admin,
		client.Pasaran_id, client.Pasaran_name, client.Pasaran_url,
		client.Pasaran_diundi, client.Pasaran_jamjadwal, client.Pasaran_status, client.Sdata, client.Pasaran_display)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_pasaran("")
	return c.JSON(result)
}
func Keluaranhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_keluaran)
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

	var obj entities.Model_keluaran
	var arraobj []entities.Model_keluaran
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Field_keluaran_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		keluaran_id, _ := jsonparser.GetInt(value, "keluaran_id")
		keluaran_tanggal, _ := jsonparser.GetString(value, "keluaran_tanggal")
		keluaran_periode, _ := jsonparser.GetString(value, "keluaran_periode")
		keluaran_nomor, _ := jsonparser.GetString(value, "keluaran_nomor")

		obj.Keluaran_id = int(keluaran_id)
		obj.Keluaran_tanggal = keluaran_tanggal
		obj.Keluaran_periode = keluaran_periode
		obj.Keluaran_nomor = keluaran_nomor
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_keluaran(client.Pasaran_id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Field_keluaran_redis+"_"+client.Pasaran_id, result, 180*time.Minute)
		log.Println("KELUARAN MYSQL")
		return c.JSON(result)
	} else {
		log.Println("KELUARAN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Keluaransave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_keluaransave)
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

	result, err := models.Save_keluaran(
		client_admin,
		client.Pasaran_id, client.Keluaran_tanggal, client.Keluaran_nomor)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_pasaran(client.Pasaran_id)
	return c.JSON(result)
}
func Keluarandelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_keluarandelete)
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

	result, err := models.Delete_keluaran(client_admin, client.Pasaran_id, client.Keluaran_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_pasaran(client.Pasaran_id)
	return c.JSON(result)
}
func Prediksihome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_keluaran)
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

	var obj entities.Model_prediksi
	var arraobj []entities.Model_prediksi
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Field_prediksi_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		prediksi_id, _ := jsonparser.GetInt(value, "prediksi_id")
		prediksi_tanggal, _ := jsonparser.GetString(value, "prediksi_tanggal")
		prediksi_bbfs, _ := jsonparser.GetString(value, "prediksi_bbfs")
		prediksi_nomor, _ := jsonparser.GetString(value, "prediksi_nomor")

		obj.Prediksi_id = int(prediksi_id)
		obj.Prediksi_tanggal = prediksi_tanggal
		obj.Prediksi_bbfs = prediksi_bbfs
		obj.Prediksi_nomor = prediksi_nomor
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_prediksi(client.Pasaran_id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Field_prediksi_redis+"_"+client.Pasaran_id, result, 180*time.Minute)
		log.Println("PREDIKSI MYSQL")
		return c.JSON(result)
	} else {
		log.Println("PREDIKSI CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Prediksisave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_prediksisave)
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

	result, err := models.Save_prediksi(
		client_admin,
		client.Pasaran_id, client.Prediksi_tanggal, client.Prediksi_bbfs, client.Prediksi_nomor)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_pasaran(client.Pasaran_id)
	return c.JSON(result)
}
func Prediksidelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_prediksidelete)
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

	result, err := models.Delete_prediksi(client_admin, client.Pasaran_id, client.Prediksi_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_pasaran(client.Pasaran_id)
	return c.JSON(result)
}
func _deleteredis_pasaran(pasaran_id string) {
	val_master := helpers.DeleteRedis(Field_home_redis)
	log.Printf("Redis Delete BACKEND PASARAN : %d", val_master)

	val_pasaran := helpers.DeleteRedis(Field_home_redis)
	val_keluaran := helpers.DeleteRedis(Field_keluaran_redis + "_" + pasaran_id)
	val_prediksi := helpers.DeleteRedis(Field_prediksi_redis + "_" + pasaran_id)
	log.Printf("Redis Delete BACKEND PASARAN : %d", val_pasaran)
	log.Printf("Redis Delete BACKEND KELUARAN : %d", val_keluaran)
	log.Printf("Redis Delete BACKEND PREDIKSI : %d", val_prediksi)

	//CLIENT
	val_client_pasaran := helpers.DeleteRedis(Field_client_pasaran_redis)
	log.Printf("Redis Delete CLIENT PASARAN : %d", val_client_pasaran)

	val_client_keluaran := helpers.DeleteRedis(Field_client_keluaran_redis + "_" + pasaran_id)
	log.Printf("Redis Delete CLIENT KELUARAN : %d", val_client_keluaran)
}
