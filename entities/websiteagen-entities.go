package entities

type Model_websiteagen struct {
	Websiteagen_id     int    `json:"websiteagen_id"`
	Websiteagen_name   string `json:"websiteagen_name"`
	Websiteagen_status string `json:"websiteagen_status"`
	Websiteagen_create string `json:"websiteagen_create"`
	Websiteagen_update string `json:"websiteagen_update"`
}

type Controller_websiteagen struct {
	Websiteagen_search string `json:"websiteagen_search"`
	Websiteagen_page   int    `json:"websiteagen_page"`
}
type Controller_websiteagensave struct {
	Page                 string `json:"page" validate:"required"`
	Sdata                string `json:"sdata" validate:"required"`
	Websiteagen_search   string `json:"websiteagen_search"`
	Websiteagen_page     int    `json:"websiteagen_page"`
	Websiteagen_idrecord int    `json:"websiteagen_id"`
	Websiteagen_name     string `json:"websiteagen_name" validate:"required"`
	Websiteagen_status   string `json:"websiteagen_status" validate:"required"`
}
