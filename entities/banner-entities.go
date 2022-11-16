package entities

type Model_banner struct {
	Banner_id         int    `json:"banner_id"`
	Banner_name       string `json:"banner_name"`
	Banner_url        string `json:"banner_url"`
	Banner_urlwebsite string `json:"banner_urlwebsite"`
	Banner_posisi     string `json:"banner_posisi"`
	Banner_device     string `json:"banner_device"`
	Banner_display    int    `json:"banner_display"`
	Banner_status     string `json:"banner_status"`
	Banner_create     string `json:"banner_create"`
	Banner_update     string `json:"banner_update"`
}
type Controller_bannersave struct {
	Sdata             string `json:"sdata" validate:"required"`
	Page              string `json:"page" validate:"required"`
	Banner_id         int    `json:"banner_id"`
	Banner_name       string `json:"banner_name" validate:"required"`
	Banner_url        string `json:"banner_url" validate:"required"`
	Banner_urlwebsite string `json:"banner_urlwebsite" validate:"required"`
	Banner_posisi     string `json:"banner_posisi" validate:"required"`
	Banner_device     string `json:"banner_device" validate:"required"`
	Banner_display    int    `json:"banner_display" `
	Banner_status     string `json:"banner_status" validate:"required"`
}
