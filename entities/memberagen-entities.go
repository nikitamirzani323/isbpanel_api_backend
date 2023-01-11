package entities

type Model_member struct {
	Member_phone  string `json:"member_phone"`
	Member_name   string `json:"member_name"`
	Member_create string `json:"member_create"`
	Member_update string `json:"member_update"`
}
type Model_memberagen struct {
	Memberagen_id        int    `json:"memberagen_id"`
	Memberagen_idwebagen int    `json:"memberagen_idwebagen"`
	Memberagen_username  string `json:"memberagen_username"`
	Memberagen_create    string `json:"memberagen_create"`
	Memberagen_update    string `json:"memberagen_update"`
}
type Controller_membersave struct {
	Sdata        string `json:"sdata" validate:"required"`
	Page         string `json:"page" validate:"required"`
	Member_phone string `json:"member_phone"`
	Member_name  string `json:"member_name" validate:"required"`
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
