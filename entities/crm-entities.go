package entities

import "encoding/json"

type Model_crm struct {
	Crm_id        int         `json:"crm_id"`
	Crm_phone     string      `json:"crm_phone"`
	Crm_name      string      `json:"crm_name"`
	Crm_pic       interface{} `json:"crm_pic"`
	Crm_totalpic  int         `json:"crm_totalpic"`
	Crm_source    string      `json:"crm_source"`
	Crm_status    string      `json:"crm_status"`
	Crm_statuscss string      `json:"crm_statuscss"`
	Crm_create    string      `json:"crm_create"`
	Crm_update    string      `json:"crm_update"`
}
type Model_crmsales_simple struct {
	Crmsales_idcrmsales   int     `json:"crmsales_idcrmsales"`
	Crmsales_username     string  `json:"crmsales_username"`
	Crmsales_nameemployee string  `json:"crmsales_nameemployee"`
	Crmsales_status_utama string  `json:"crmsales_status_utama"`
	Crmsales_status       string  `json:"crmsales_status"`
	Crmsales_note         string  `json:"crmsales_note"`
	Crmsales_nmwebagen    string  `json:"crmsales_nmwebagen"`
	Crmsales_idwebagen    string  `json:"crmsales_idwebagen"`
	Crmsales_deposit      float32 `json:"crmsales_deposit"`
}
type Model_crmsales struct {
	Crmsales_id           int    `json:"crmsales_id"`
	Crmsales_phone        string `json:"crmsales_phone"`
	Crmsales_namamember   string `json:"crmsales_namamember"`
	Crmsales_username     string `json:"crmsales_username"`
	Crmsales_nameemployee string `json:"crmsales_nameemployee"`
	Crmsales_create       string `json:"crmsales_create"`
	Crmsales_update       string `json:"crmsales_update"`
}
type Model_crmdeposit struct {
	Crmsdeposit_nmwebagen  string  `json:"crmdeposit_nmwebagen"`
	Crmsdeposit_deposit    float32 `json:"crmdeposit_deposit"`
	Crmsdeposit_iduseragen string  `json:"crmdeposit_iduseragen"`
	Crmsdeposit_create     string  `json:"crmdeposit_create"`
}
type Model_crmmemberlistdeposit struct {
	Crmsdeposit_phone      string  `json:"crmdeposit_phone"`
	Crmsdeposit_nama       string  `json:"crmdeposit_nama"`
	Crmsdeposit_source     string  `json:"crmdeposit_source"`
	Crmsdeposit_nmwebagen  string  `json:"crmdeposit_nmwebagen"`
	Crmsdeposit_deposit    float32 `json:"crmdeposit_deposit"`
	Crmsdeposit_iduseragen string  `json:"crmdeposit_iduseragen"`
	Crmsdeposit_update     string  `json:"crmdeposit_update"`
}
type Model_crmmemberlistnoanswer struct {
	Crmnoanswer_phone  string `json:"crmnoanswer_phone"`
	Crmnoanswer_nama   string `json:"crmnoanswer_nama"`
	Crmnoanswer_source string `json:"crmnoanswer_source"`
	Crmnoanswer_tipe   string `json:"crmnoanswer_tipe"`
	Crmnoanswer_note   string `json:"crmnoanswer_note"`
	Crmnoanswer_update string `json:"crmnoanswer_update"`
}
type Model_crmmemberlistinvalid struct {
	Crminvalid_phone  string `json:"crminvalid_phone"`
	Crminvalid_nama   string `json:"crminvalid_nama"`
	Crminvalid_source string `json:"crminvalid_source"`
	Crminvalid_update string `json:"crminvalid_update"`
}
type Model_crmisbtv struct {
	Crmisbtv_username  string `json:"crmisbtv_username"`
	Crmisbtv_name      string `json:"crmisbtv_name"`
	Crmisbtv_coderef   string `json:"crmisbtv_coderef"`
	Crmisbtv_point     int    `json:"crmisbtv_point"`
	Crmisbtv_status    string `json:"crmisbtv_status"`
	Crmisbtv_lastlogin string `json:"crmisbtv_lastlogin"`
	Crmisbtv_create    string `json:"crmisbtv_create"`
	Crmisbtv_update    string `json:"crmisbtv_update"`
}
type Model_crmduniafilm struct {
	Crmduniafilm_username string `json:"crmduniafilm_username"`
	Crmduniafilm_name     string `json:"crmduniafilm_name"`
}

type Controller_crm struct {
	Crm_status string `json:"crm_status"`
	Crm_search string `json:"crm_search"`
	Crm_page   int    `json:"crm_page"`
}
type Controller_crmsales struct {
	Crmsales_phone  string `json:"crmsales_phone"`
	Crmsales_status string `json:"crmsales_status"`
}
type Controller_crmdeposit struct {
	Crmsales_idcrmsales int `json:"crmsales_idcrmsales" validate:"required"`
}
type Controller_crmisbtv struct {
	Crmisbtv_search string `json:"crmisbtv_search"`
	Crmisbtv_page   int    `json:"crmisbtv_page"`
}
type Controller_crmsave struct {
	Page       string `json:"page" validate:"required"`
	Sdata      string `json:"sdata" validate:"required"`
	Crm_page   int    `json:"crm_page"`
	Crm_id     int    `json:"crm_id"`
	Crm_phone  string `json:"crm_phone" validate:"required"`
	Crm_name   string `json:"crm_name" validate:"required"`
	Crm_status string `json:"crm_status" validate:"required"`
}
type Controller_crmstatussave struct {
	Page       string `json:"page" validate:"required"`
	Crm_page   int    `json:"crm_page"`
	Crm_id     int    `json:"crm_id"`
	Crm_status string `json:"crm_status" validate:"required"`
}
type Controller_crmsalessave struct {
	Page              string `json:"page" validate:"required"`
	Search            string `json:"search" `
	Crm_page          int    `json:"crm_page"`
	Crmsales_phone    string `json:"crmsales_phone" validate:"required"`
	Crmsales_username string `json:"crmsales_username" validate:"required"`
}
type Controller_crmsavesource struct {
	Page       string          `json:"page" validate:"required"`
	Sdata      string          `json:"sdata" validate:"required"`
	Crm_page   int             `json:"crm_page"`
	Crm_source string          `json:"crm_source" `
	Crm_data   json.RawMessage `json:"crm_data" validate:"required"`
}
type Controller_crmsavemaintenance struct {
	Page     string          `json:"page" validate:"required"`
	Sdata    string          `json:"sdata" validate:"required"`
	Crm_page int             `json:"crm_page"`
	Crm_data json.RawMessage `json:"crm_data" validate:"required"`
}
type Controller_crmsalesdelete struct {
	Page           string `json:"page" validate:"required"`
	Search         string `json:"search" `
	Crm_page       int    `json:"crm_page"`
	Crmsales_id    int    `json:"crmsales_id" validate:"required"`
	Crmsales_phone string `json:"crmsales_phone" validate:"required"`
}
