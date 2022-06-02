package entities

import "encoding/json"

type Model_album struct {
	Album_id             int    `json:"album_id"`
	Album_idcloud        string `json:"album_idcloud"`
	Album_name           string `json:"album_name"`
	Album_signed         string `json:"album_signed"`
	Album_varian         string `json:"album_varian"`
	Album_movieid        int    `json:"album_movieid"`
	Album_movie          string `json:"album_movie"`
	Album_moviestatus    string `json:"album_moviestatus"`
	Album_moviestatuscss string `json:"album_moviestatuscss"`
	Album_createdate     string `json:"album_createdate"`
}

type Controller_album struct {
	Album_page int `json:"album_page"`
}
type Controller_albumcloudflare struct {
	Album_page      int `json:"album_page"`
	Album_pagecloud int `json:"album_pagecloud"`
}
type Controller_cloudflaresave struct {
	Page       string          `json:"page" validate:"required"`
	Sdata      string          `json:"sdata" validate:"required"`
	Album_page int             `json:"album_page" `
	Album_data json.RawMessage `json:"album_data" validate:"required"`
}
