package entities

type Model_prediksislot struct {
	Prediksislot_id         int    `json:"prediksislot_id"`
	Prediksislot_nmprovider string `json:"prediksislot_nmprovider"`
	Prediksislot_name       string `json:"prediksislot_name"`
	Prediksislot_prediksi   int    `json:"prediksislot_prediksi"`
	Prediksislot_image      string `json:"prediksislot_image"`
	Prediksislot_status     string `json:"prediksislot_status"`
	Prediksislot_create     string `json:"prediksislot_create"`
	Prediksislot_update     string `json:"prediksislot_update"`
}
type Controller_prediksislothome struct {
	Providerslot_id int `json:"providerslot_id" `
}
type Controller_prediksislotsave struct {
	Page                  string `json:"page" validate:"required"`
	Sdata                 string `json:"sdata" validate:"required"`
	Prediksislot_id       int    `json:"prediksislot_id"`
	Providerslot_id       int    `json:"providerslot_id" validate:"required"`
	Providerslot_slug     string `json:"providerslot_slug" validate:"required"`
	Prediksislot_name     string `json:"prediksislot_name" validate:"required"`
	Prediksislot_prediksi int    `json:"prediksislot_prediksi"`
	Prediksislot_image    string `json:"prediksislot_image" validate:"required"`
	Prediksislot_status   string `json:"prediksislot_status" validate:"required"`
}
type Controller_prediksislotdelete struct {
	Page              string `json:"page" validate:"required"`
	Prediksislot_id   int    `json:"prediksislot_id" validate:"required"`
	Providerslot_id   int    `json:"providerslot_id" validate:"required"`
	Providerslot_slug string `json:"providerslot_slug" validate:"required"`
}
type Controller_prediksislotgenerator struct {
	Page              string `json:"page" validate:"required"`
	Providerslot_id   int    `json:"providerslot_id" validate:"required"`
	Providerslot_slug string `json:"providerslot_slug" validate:"required"`
}
