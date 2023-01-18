package entities

import "encoding/json"

type Model_member struct {
	Member_phone  string      `json:"member_phone"`
	Member_name   string      `json:"member_name"`
	Member_agen   interface{} `json:"member_agen"`
	Member_create string      `json:"member_create"`
	Member_update string      `json:"member_update"`
}
type Model_memberagen struct {
	Memberagen_idwebagen int    `json:"memberagen_idwebagen"`
	Memberagen_website   string `json:"memberagen_website"`
	Memberagen_username  string `json:"memberagen_username"`
}
type Model_memberagenselect struct {
	Memberagen_id       int    `json:"memberagen_id"`
	Memberagen_website  string `json:"memberagen_website"`
	Memberagen_username string `json:"memberagen_username"`
	Memberagen_phone    string `json:"memberagen_phone"`
	Memberagen_name     string `json:"memberagen_name"`
}

type Controller_membersave struct {
	Sdata           string          `json:"sdata" validate:"required"`
	Page            string          `json:"page" validate:"required"`
	Member_phone    string          `json:"member_phone" validate:"required"`
	Member_name     string          `json:"member_name" validate:"required"`
	Member_listagen json.RawMessage `json:"member_listagen" validate:"required"`
}
type Controller_memberagen struct {
	Memberagen_phone string `json:"memberagen_phone" validate:"required"`
}
type Controller_memberagenselect struct {
	Memberagen_idwebagen int `json:"memberagen_idwebagen" validate:"required"`
}
