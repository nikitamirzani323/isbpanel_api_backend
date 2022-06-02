package entities

type Model_tafsirmimpi struct {
	Tafsirmimpi_id        int    `json:"tafsirmimpi_id"`
	Tafsirmimpi_mimpi     string `json:"tafsirmimpi_mimpi"`
	Tafsirmimpi_artimimpi string `json:"tafsirmimpi_artimimpi"`
	Tafsirmimpi_angka2d   string `json:"tafsirmimpi_angka2d"`
	Tafsirmimpi_angka3d   string `json:"tafsirmimpi_angka3d"`
	Tafsirmimpi_angka4d   string `json:"tafsirmimpi_angka4d"`
	Tafsirmimpi_status    string `json:"tafsirmimpi_status"`
	Tafsirmimpi_statuscss string `json:"tafsirmimpi_statuscss"`
	Tafsirmimpi_create    string `json:"tafsirmimpi_create"`
	Tafsirmimpi_update    string `json:"tafsirmimpi_update"`
}
type Controller_tafsirmimpi struct {
	Tafsirmimpi_search string `json:"tafsirmimpi_search"`
}
type Controller_tafsirmimpisave struct {
	Sdata                 string `json:"sdata" validate:"required"`
	Page                  string `json:"page" validate:"required"`
	Tafsirmimpi_id        int    `json:"tafsirmimpi_id"`
	Tafsirmimpi_mimpi     string `json:"tafsirmimpi_mimpi" validate:"required"`
	Tafsirmimpi_artimimpi string `json:"tafsirmimpi_artimimpi" validate:"required"`
	Tafsirmimpi_angka2d   string `json:"tafsirmimpi_angka2d" validate:"required"`
	Tafsirmimpi_angka3d   string `json:"tafsirmimpi_angka3d" validate:"required"`
	Tafsirmimpi_angka4d   string `json:"tafsirmimpi_angka4d" validate:"required"`
	Tafsirmimpi_status    string `json:"tafsirmimpi_status"`
}
