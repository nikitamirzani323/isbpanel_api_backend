package controllers

import (
	"fmt"
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

const Fieldleague_home_redis = "LISTLEAGUE_BACKEND_ISBPANEL"
const Fieldleague_frontend_redis = "LISTLEAGUE_FRONTEND_ISBPANEL"

func Leaguehome(c *fiber.Ctx) error {
	var obj entities.Model_league
	var arraobj []entities.Model_league
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldleague_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		league_id, _ := jsonparser.GetInt(value, "league_id")
		league_name, _ := jsonparser.GetString(value, "league_name")
		league_image, _ := jsonparser.GetString(value, "league_image")
		league_status, _ := jsonparser.GetString(value, "league_status")
		league_create, _ := jsonparser.GetString(value, "league_create")
		league_update, _ := jsonparser.GetString(value, "league_update")

		obj.League_id = int(league_id)
		obj.League_name = league_name
		obj.League_image = league_image
		obj.League_status = league_status
		obj.League_create = league_create
		obj.League_update = league_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_leagueHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldleague_home_redis, result, 60*time.Minute)
		fmt.Println("LEAGUE MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("LEAGUE CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func LeagueSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_leaguesave)
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

	result, err := models.Save_league(
		client_admin,
		client.League_name, client.League_image, client.League_status, client.Sdata, client.League_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_league()
	return c.JSON(result)
}
func _deleteredis_league() {
	val_master := helpers.DeleteRedis(Fieldleague_home_redis)
	log.Printf("Redis Delete BACKEND LEAGUE : %d", val_master)

	val_client := helpers.DeleteRedis(Fieldleague_frontend_redis)
	log.Printf("Redis Delete FRONTEND DOMAIN : %d", val_client)

}
