package entities

import "encoding/json"

type Model_movie struct {
	Movie_id        int         `json:"movie_id"`
	Movie_date      string      `json:"movie_date"`
	Movie_type      string      `json:"movie_type"`
	Movie_title     string      `json:"movie_title"`
	Movie_label     string      `json:"movie_label"`
	Movie_slug      string      `json:"movie_slug"`
	Movie_descp     string      `json:"movie_descp"`
	Movie_imgcdn    string      `json:"movie_imgcdn"`
	Movie_thumbnail string      `json:"movie_thumbnail"`
	Movie_year      int         `json:"movie_year"`
	Movie_rating    float32     `json:"movie_rating"`
	Movie_imdb      float32     `json:"movie_imdb"`
	Movie_view      int         `json:"movie_view"`
	Movie_comment   int         `json:"movie_comment"`
	Movie_genre     interface{} `json:"movie_genre"`
	Movie_source    interface{} `json:"movie_source"`
	Movie_status    string      `json:"movie_status"`
	Movie_statuscss string      `json:"movie_statuscss"`
	Movie_create    string      `json:"movie_create"`
	Movie_update    string      `json:"movie_update"`
}
type Model_movienotcdn struct {
	Movie_id    int    `json:"movie_id"`
	Movie_title string `json:"movie_title"`
}
type Model_moviegenre struct {
	Moviegenre_id   int    `json:"moviegenre_id"`
	Moviegenre_name string `json:"moviegenre_name"`
}
type Model_moviesource struct {
	Moviesource_id     int    `json:"moviesource_id"`
	Moviesource_stream string `json:"moviesource_stream"`
	Moviesource_url    string `json:"moviesource_url"`
}
type Model_movieseries struct {
	Movie_id        int         `json:"movie_id"`
	Movie_date      string      `json:"movie_date"`
	Movie_type      string      `json:"movie_type"`
	Movie_title     string      `json:"movie_title"`
	Movie_label     string      `json:"movie_label"`
	Movie_slug      string      `json:"movie_slug"`
	Movie_descp     string      `json:"movie_descp"`
	Movie_imgcdn    string      `json:"movie_imgcdn"`
	Movie_thumbnail string      `json:"movie_thumbnail"`
	Movie_year      int         `json:"movie_year"`
	Movie_rating    float32     `json:"movie_rating"`
	Movie_imdb      float32     `json:"movie_imdb"`
	Movie_view      int         `json:"movie_view"`
	Movie_comment   int         `json:"movie_comment"`
	Movie_genre     interface{} `json:"movie_genre"`
	Movie_season    interface{} `json:"movie_season"`
	Movie_status    string      `json:"movie_status"`
	Movie_statuscss string      `json:"movie_statuscss"`
	Movie_create    string      `json:"movie_create"`
	Movie_update    string      `json:"movie_update"`
}
type Model_movieseason struct {
	Movieseason_id           int    `json:"movieseason_id"`
	Movieseason_name         string `json:"movieseason_name"`
	Movieseason_episodetotal int    `json:"movieseason_episodetotal"`
	Movieseason_display      int    `json:"movieseason_display"`
}
type Model_movieepisode struct {
	Movieepisode_id       int         `json:"movieepisode_id"`
	Movieepisode_seasonid int         `json:"movieepisode_seasonid"`
	Movieepisode_name     string      `json:"movieepisode_name"`
	Movieepisode_display  int         `json:"movieepisode_display"`
	Movieepisode_source   interface{} `json:"movieepisode_source"`
}
type Model_genre struct {
	Genre_id      int    `json:"genre_id"`
	Genre_name    string `json:"genre_name"`
	Genre_display int    `json:"genre_display"`
	Genre_create  string `json:"genre_create"`
	Genre_update  string `json:"genre_update"`
}

type Model_minimovie struct {
	Movie_id    int    `json:"movie_id"`
	Movie_type  string `json:"movie_type"`
	Movie_title string `json:"movie_title"`
}
type Controller_movie struct {
	Movie_search string `json:"movie_search"`
	Movie_page   int    `json:"movie_page"`
}
type Controller_moviemini struct {
	Movie_search string `json:"movie_search"`
}
type Controller_moviesave struct {
	Page           string          `json:"page" validate:"required"`
	Sdata          string          `json:"sdata" validate:"required"`
	Movie_page     int             `json:"movie_page" `
	Movie_id       int             `json:"movie_id"`
	Movie_name     string          `json:"movie_name" validate:"required"`
	Movie_label    string          `json:"movie_label" validate:"required"`
	Movie_slug     string          `json:"movie_slug" validate:"required"`
	Movie_tipe     string          `json:"movie_tipe" validate:"required"`
	Movie_descp    string          `json:"movie_descp" validate:"required"`
	Movie_urlmovie string          `json:"movie_urlmovie" validate:"required"`
	Movie_year     int             `json:"movie_year" validate:"required"`
	Movie_imdb     float32         `json:"movie_imdb" validate:"required"`
	Movie_status   int             `json:"movie_status"`
	Movie_gender   json.RawMessage `json:"movie_gender" validate:"required"`
	Movie_source   json.RawMessage `json:"movie_source" `
}
type Controller_moviedelete struct {
	Page       string `json:"page" validate:"required"`
	Movie_page int    `json:"movie_page" `
	Movie_id   int    `json:"movie_id" validate:"required"`
}
type Controller_movieseriessave struct {
	Page           string          `json:"page" validate:"required"`
	Sdata          string          `json:"sdata" validate:"required"`
	Movie_page     int             `json:"movie_page" `
	Movie_id       int             `json:"movie_id"`
	Movie_name     string          `json:"movie_name" validate:"required"`
	Movie_label    string          `json:"movie_label" validate:"required"`
	Movie_slug     string          `json:"movie_slug" validate:"required"`
	Movie_tipe     string          `json:"movie_tipe" validate:"required"`
	Movie_descp    string          `json:"movie_descp" validate:"required"`
	Movie_urlmovie string          `json:"movie_urlmovie" validate:"required"`
	Movie_year     int             `json:"movie_year" validate:"required"`
	Movie_imdb     float32         `json:"movie_imdb" validate:"required"`
	Movie_status   int             `json:"movie_status"`
	Movie_gender   json.RawMessage `json:"movie_gender" validate:"required"`
}
type Controller_movieseason struct {
	Movie_id int `json:"movie_id" validate:"required"`
}
type Controller_movieseasonsave struct {
	Page                string `json:"page" validate:"required"`
	Sdata               string `json:"sdata" validate:"required"`
	Movie_page          int    `json:"movie_page" `
	Movie_id            int    `json:"movie_id" validate:"required"`
	Movieseason_id      int    `json:"movieseason_id"`
	Movieseason_name    string `json:"movieseason_name" validate:"required"`
	Movieseason_display int    `json:"movieseason_display" validate:"required"`
}
type Controller_movieseasondelete struct {
	Page           string `json:"page" validate:"required"`
	Movie_page     int    `json:"movie_page" `
	Movie_id       int    `json:"movie_id" validate:"required"`
	Movieseason_id int    `json:"movieseason_id" validate:"required"`
}
type Controller_movieepisode struct {
	Movie_id  int `json:"movie_id" validate:"required"`
	Season_id int `json:"season_id" validate:"required"`
}
type Controller_movieepisodesave struct {
	Page                 string `json:"page" validate:"required"`
	Sdata                string `json:"sdata" validate:"required"`
	Movie_page           int    `json:"movie_page" `
	Movie_id             int    `json:"movie_id" validate:"required"`
	Movieseason_id       int    `json:"movieseason_id" alidate:"required"`
	Movieepisode_id      int    `json:"movieepisode_id"`
	Movieepisode_name    string `json:"movieepisode_name" validate:"required"`
	Movieepisode_source  string `json:"movieepisode_source" validate:"required"`
	Movieepisode_display int    `json:"movieepisode_display" validate:"required"`
}
type Controller_movieepisodedelete struct {
	Page       string `json:"page" validate:"required"`
	Movie_page int    `json:"movie_page" `
	Movie_id   int    `json:"movie_id" validate:"required"`
	Season_id  int    `json:"season_id" validate:"required"`
	Episode_id int    `json:"episode_id" validate:"required"`
}
type Controller_cloudflaremovieupload struct {
	Page      string `json:"page" validate:"required"`
	Sdata     string `json:"sdata" validate:"required"`
	Movie_raw string `json:"movie_raw" validate:"required"`
}
type Controller_cloudflaremovieupdate struct {
	Page       string `json:"page" validate:"required"`
	Sdata      string `json:"sdata" validate:"required"`
	Movie_id   string `json:"movie_id" validate:"required"`
	Movie_tipe string `json:"movie_tipe" validate:"required"`
	Album_id   int    `json:"album_id" `
}
type Controller_cloudflaremoviedelete struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	Cloudflare_id string `json:"cloudflare_id" validate:"required"`
	Movie_id      int    `json:"movie_id" `
	Album_id      int    `json:"album_id" validate:"required"`
}
type Controller_genresave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	Genre_id      int    `json:"genre_id"`
	Genre_name    string `json:"genre_name" validate:"required"`
	Genre_display int    `json:"genre_display" validate:"required"`
}
type Controller_genredelete struct {
	Page     string `json:"page" validate:"required"`
	Sdata    string `json:"sdata" validate:"required"`
	Genre_id int    `json:"genre_id" validate:"required"`
}
