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

const Fieldnews_home_redis = "LISTNEWS_BACKEND_ISBPANEL"
const Fieldcategory_home_redis = "LISTCATEGORY_BACKEND_ISBPANEL"
const Fieldnews_client_home_redis = "LISTNEWS_FRONTEND_ISBPANEL"
const Fieldnewsmovie_client_home_redis = "LISTNEWSMOVIES_FRONTEND_ISBPANEL"

func Newshome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_news)
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
	if client.News_search != "" {
		val_news := helpers.DeleteRedis(Fieldnews_home_redis + "_" + strconv.Itoa(client.News_page) + "_" + client.News_search)
		log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	}
	var obj entities.Model_news
	var arraobj []entities.Model_news
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldnews_home_redis + "_" + strconv.Itoa(client.News_page) + "_" + client.News_search)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		news_id, _ := jsonparser.GetInt(value, "news_id")
		news_title, _ := jsonparser.GetString(value, "news_title")
		news_idcategory, _ := jsonparser.GetInt(value, "news_idcategory")
		news_category, _ := jsonparser.GetString(value, "news_category")
		news_descp, _ := jsonparser.GetString(value, "news_descp")
		news_url, _ := jsonparser.GetString(value, "news_url")
		news_image, _ := jsonparser.GetString(value, "news_image")
		news_create, _ := jsonparser.GetString(value, "news_create")
		news_update, _ := jsonparser.GetString(value, "news_update")

		obj.News_id = int(news_id)
		obj.News_idcategory = int(news_idcategory)
		obj.News_category = news_category
		obj.News_title = news_title
		obj.News_descp = news_descp
		obj.News_url = news_url
		obj.News_image = news_image
		obj.News_create = news_create
		obj.News_update = news_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_newsHome(client.News_search, client.News_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldnews_home_redis+"_"+strconv.Itoa(client.News_page)+"_"+client.News_search, result, 10*time.Minute)
		log.Println("NEWS MYSQL")
		return c.JSON(result)
	} else {
		log.Println("NEWS CACHE")
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
func Newssave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_newssave)
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

	result, err := models.Save_news(
		client_admin, client.Sdata,
		client.News_title, client.News_descp, client.News_url, client.News_image, client.News_idrecord, client.News_category)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_news := helpers.DeleteRedis(Fieldnews_home_redis + "_")
	log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	val_client_news := helpers.DeleteRedis(Fieldnews_client_home_redis)
	log.Printf("Redis Delete CLIENT NEWS : %d", val_client_news)
	val_client_newsmovie := helpers.DeleteRedis(Fieldnewsmovie_client_home_redis)
	log.Printf("Redis Delete CLIENT NEWS MOVIE : %d", val_client_newsmovie)
	return c.JSON(result)
}
func Newsdelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_newsdelete)
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

	result, err := models.Delete_news(client_admin, client.News_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_news := helpers.DeleteRedis(Fieldnews_home_redis + "_")
	log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	val_client_news := helpers.DeleteRedis(Fieldnews_client_home_redis)
	log.Printf("Redis Delete CLIENT NEWS : %d", val_client_news)
	val_client_newsmovie := helpers.DeleteRedis(Fieldnewsmovie_client_home_redis)
	log.Printf("Redis Delete CLIENT NEWS MOVIE : %d", val_client_newsmovie)
	return c.JSON(result)
}
func Categoryhome(c *fiber.Ctx) error {
	var obj entities.Model_category
	var arraobj []entities.Model_category
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcategory_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		category_id, _ := jsonparser.GetInt(value, "category_id")
		category_name, _ := jsonparser.GetString(value, "category_name")
		category_totalnews, _ := jsonparser.GetInt(value, "category_totalnews")
		category_display, _ := jsonparser.GetInt(value, "category_display")
		category_status, _ := jsonparser.GetString(value, "category_status")
		category_statuscss, _ := jsonparser.GetString(value, "category_statuscss")
		category_create, _ := jsonparser.GetString(value, "category_create")
		category_update, _ := jsonparser.GetString(value, "category_update")

		obj.Category_id = int(category_id)
		obj.Category_name = category_name
		obj.Category_totalnews = int(category_totalnews)
		obj.Category_display = int(category_display)
		obj.Category_status = category_status
		obj.Category_statuscss = category_statuscss
		obj.Category_create = category_create
		obj.Category_update = category_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_category()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcategory_home_redis, result, 60*time.Minute)
		log.Println("CATEGORY NEWS MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CATEGORY NEWS CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Categorysave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_categorysave)
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

	result, err := models.Save_category(
		client_admin,
		client.Category_name, client.Category_status, client.Sdata, client.Category_idrecord, client.Category_display)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_category := helpers.DeleteRedis(Fieldcategory_home_redis)
	log.Printf("Redis Delete BACKEND NEWS CATEGORY : %d", val_category)
	return c.JSON(result)
}
func Categorydelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_categorydelete)
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

	result, err := models.Delete_category(client_admin, client.Category_idrecord)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	val_category := helpers.DeleteRedis(Fieldcategory_home_redis)
	log.Printf("Redis Delete BACKEND NEWS CATEGORY : %d", val_category)
	return c.JSON(result)
}
