package entities

type Model_event struct {
	Event_id         int    `json:"event_id"`
	Event_idwebagen  int    `json:"event_idwebagen"`
	Event_name       string `json:"event_name"`
	Event_startevent string `json:"event_startevent"`
	Event_endevent   string `json:"event_endevent"`
	Event_create     string `json:"event_create"`
	Event_update     string `json:"event_update"`
}
type Controller_eventsave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Event_id         int    `json:"event_id"`
	Event_idwebagen  int    `json:"event_idwebagen" validate:"required"`
	Event_name       string `json:"event_name" validate:"required"`
	Event_startevent string `json:"event_startevent" validate:"required"`
	Event_endevent   string `json:"event_endevent" validate:"required"`
}
