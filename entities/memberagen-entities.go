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
	Memberagen_id       int    `json:"memberagen_id"`
	Memberagen_website  string `json:"memberagen_website"`
	Memberagen_username string `json:"memberagen_username"`
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
type Controller_memberagensave struct {
	Sdata                string `json:"sdata" validate:"required"`
	Page                 string `json:"page" validate:"required"`
	Memberagen_id        int    `json:"memberagen_id"`
	Memberagen_idwebagen int    `json:"memberagen_idwebagen" validate:"required"`
	Memberagen_phone     string `json:"memberagen_phone" validate:"required"`
	Memberagen_username  string `json:"memberagen_username" validate:"required"`
}
