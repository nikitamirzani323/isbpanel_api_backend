package entities

type Model_event struct {
	Event_id         int    `json:"event_id"`
	Event_idwebagen  int    `json:"event_idwebagen"`
	Event_nmwebagen  string `json:"event_nmwebagen"`
	Event_name       string `json:"event_name"`
	Event_startevent string `json:"event_startevent"`
	Event_endevent   string `json:"event_endevent"`
	Event_mindeposit int    `json:"event_mindeposit"`
	Event_create     string `json:"event_create"`
	Event_update     string `json:"event_update"`
}
type Model_eventdetail struct {
	Eventdetail_iddetail int    `json:"eventdetail_id"`
	Eventdetail_voucher  string `json:"eventdetail_voucher"`
	Eventdetail_deposit  int    `json:"eventdetail_deposit"`
	Eventdetail_phone    string `json:"eventdetail_phone"`
	Eventdetail_username string `json:"eventdetail_username"`
	Eventdetail_create   string `json:"eventdetail_create"`
	Eventdetail_update   string `json:"eventdetail_update"`
}
type Model_eventdetailgroup struct {
	Eventdetailgroup_idmember int    `json:"eventdetailgroup_idmember"`
	Eventdetailgroup_deposit  int    `json:"eventdetailgroup_deposit"`
	Eventdetailgroup_phone    string `json:"eventdetailgroup_phone"`
	Eventdetailgroup_username string `json:"eventdetailgroup_username"`
}
type Controller_eventdetail struct {
	Sdata              string `json:"sdata" validate:"required"`
	Page               string `json:"page" validate:"required"`
	Event_id           int    `json:"event_id"`
	Event_idmemberagen int    `json:"event_idmemberagen"`
}
type Controller_eventsave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Event_id         int    `json:"event_id"`
	Event_idwebagen  int    `json:"event_idwebagen" validate:"required"`
	Event_name       string `json:"event_name" validate:"required"`
	Event_startevent string `json:"event_startevent" validate:"required"`
	Event_endevent   string `json:"event_endevent" validate:"required"`
	Event_mindeposit int    `json:"event_mindeposit" validate:"required"`
}
type Controller_eventdetailsave struct {
	Sdata                    string `json:"sdata" validate:"required"`
	Page                     string `json:"page" validate:"required"`
	Eventdetail_id           int    `json:"eventdetail_id"`
	Eventdetail_idevent      int    `json:"eventdetail_idevent"`
	Eventdetail_idmemberagen int    `json:"eventdetail_idmemberagen" validate:"required"`
	Eventdetail_deposit      int    `json:"eventdetail_deposit" validate:"required"`
}
