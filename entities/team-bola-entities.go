package entities

type Model_Team struct {
	Team_id     int    `json:"team_id"`
	Team_name   string `json:"team_name"`
	Team_image  string `json:"team_image"`
	Team_status string `json:"team_status"`
	Team_create string `json:"team_create"`
	Team_update string `json:"team_update"`
}

type Controller_Teamsave struct {
	Page        string `json:"page" validate:"required"`
	Sdata       string `json:"sdata" validate:"required"`
	Team_id     int    `json:"team_id"`
	Team_name   string `json:"team_name" validate:"required"`
	Team_image  string `json:"team_image" validate:"required"`
	Team_status string `json:"team_status" validate:"required"`
}
