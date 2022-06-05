package entities

type Model_providerslot struct {
	Providerslot_id            int    `json:"providerslot_id"`
	Providerslot_name          string `json:"providerslot_name"`
	Providerslot_display       int    `json:"providerslot_display"`
	Providerslot_counter       int    `json:"providerslot_counter"`
	Providerslot_totalgameslot int    `json:"providerslot_totalgameslot"`
	Providerslot_image         string `json:"providerslot_image"`
	Providerslot_slug          string `json:"providerslot_slug"`
	Providerslot_title         string `json:"providerslot_title"`
	Providerslot_descp         string `json:"providerslot_descp"`
	Providerslot_status        string `json:"providerslot_status"`
	Providerslot_create        string `json:"providerslot_create"`
	Providerslot_update        string `json:"providerslot_update"`
}
type Controller_providerslotsave struct {
	Page                 string `json:"page" validate:"required"`
	Sdata                string `json:"sdata" validate:"required"`
	Providerslot_id      int    `json:"providerslot_id"`
	Providerslot_name    string `json:"providerslot_name" validate:"required"`
	Providerslot_display int    `json:"providerslot_display" validate:"required"`
	Providerslot_image   string `json:"providerslot_image" validate:"required"`
	Providerslot_slug    string `json:"providerslot_slug" validate:"required"`
	Providerslot_title   string `json:"providerslot_title" validate:"required"`
	Providerslot_descp   string `json:"providerslot_descp" validate:"required"`
	Providerslot_status  string `json:"providerslot_status" validate:"required"`
}
