package entities

type Model_league struct {
	League_id     int    `json:"league_id"`
	League_name   string `json:"league_name"`
	League_image  string `json:"league_image"`
	League_status string `json:"league_status"`
	League_create string `json:"league_create"`
	League_update string `json:"league_update"`
}

type Controller_leaguesave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	League_id     int    `json:"league_id"`
	League_name   string `json:"league_name" validate:"required"`
	League_image  string `json:"league_image" validate:"required"`
	League_status string `json:"league_status" validate:"required"`
}
