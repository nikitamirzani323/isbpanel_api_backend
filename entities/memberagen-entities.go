package entities

type Model_memberagen struct {
	Memberagen_id        int    `json:"memberagen_id"`
	Memberagen_idwebagen int    `json:"memberagen_idwebagen"`
	Memberagen_username  string `json:"memberagen_username"`
	Memberagen_name      string `json:"memberagen_name"`
	Memberagen_create    string `json:"memberagen_create"`
	Memberagen_update    string `json:"memberagen_update"`
}
type Controller_memberagensave struct {
	Sdata                string `json:"sdata" validate:"required"`
	Page                 string `json:"page" validate:"required"`
	Memberagen_id        int    `json:"memberagen_id"`
	Memberagen_idwebagen int    `json:"memberagen_idwebagen" validate:"required"`
	Memberagen_username  string `json:"memberagen_username" validate:"required"`
	Memberagen_name      string `json:"memberagen_name" validate:"required"`
}
