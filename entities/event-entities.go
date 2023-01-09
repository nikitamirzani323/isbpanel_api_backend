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
