package entities

type Model_pasaran struct {
	Pasaran_id        string `json:"pasaran_id"`
	Pasaran_name      string `json:"pasaran_name"`
	Pasaran_url       string `json:"pasaran_url"`
	Pasaran_diundi    string `json:"pasaran_diundi"`
	Pasaran_jamjadwal string `json:"pasaran_jamjadwal"`
	Pasaran_display   int    `json:"pasaran_display"`
	Pasaran_status    string `json:"pasaran_status"`
	Pasaran_statuscss string `json:"pasaran_statuscss"`
	Pasaran_keluaran  string `json:"pasaran_keluaran"`
	Pasaran_prediksi  string `json:"pasaran_prediksi"`
	Pasaran_create    string `json:"pasaran_create"`
	Pasaran_update    string `json:"pasaran_update"`
}
type Model_keluaran struct {
	Keluaran_id      int    `json:"keluaran_id"`
	Keluaran_tanggal string `json:"keluaran_tanggal"`
	Keluaran_periode string `json:"keluaran_periode"`
	Keluaran_nomor   string `json:"keluaran_nomor"`
}
type Model_prediksi struct {
	Prediksi_id      int    `json:"prediksi_id"`
	Prediksi_tanggal string `json:"prediksi_tanggal"`
	Prediksi_bbfs    string `json:"prediksi_bbfs"`
	Prediksi_nomor   string `json:"prediksi_nomor"`
}
type Controller_pasaransave struct {
	Sdata             string `json:"sdata" validate:"required"`
	Page              string `json:"page" validate:"required"`
	Pasaran_id        string `json:"pasaran_id"`
	Pasaran_name      string `json:"pasaran_name" validate:"required"`
	Pasaran_url       string `json:"pasaran_url" validate:"required"`
	Pasaran_diundi    string `json:"pasaran_diundi" validate:"required"`
	Pasaran_jamjadwal string `json:"pasaran_jamjadwal" validate:"required"`
	Pasaran_display   int    `json:"pasaran_display" validate:"required"`
	Pasaran_status    string `json:"pasaran_status" `
}
type Controller_keluaran struct {
	Page       string `json:"page" validate:"required"`
	Pasaran_id string `json:"pasaran_id" validate:"required"`
}
type Controller_keluaransave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Pasaran_id       string `json:"pasaran_id"`
	Keluaran_tanggal string `json:"keluaran_tanggal"`
	Keluaran_nomor   string `json:"keluaran_nomor"`
}
type Controller_keluarandelete struct {
	Page        string `json:"page" validate:"required"`
	Pasaran_id  string `json:"pasaran_id"`
	Keluaran_id int    `json:"keluaran_id"`
}
type Controller_prediksisave struct {
	Sdata            string `json:"sdata" validate:"required"`
	Page             string `json:"page" validate:"required"`
	Pasaran_id       string `json:"pasaran_id"`
	Prediksi_tanggal string `json:"prediksi_tanggal"`
	Prediksi_bbfs    string `json:"prediksi_bbfs"`
	Prediksi_nomor   string `json:"prediksi_nomor"`
}
type Controller_prediksidelete struct {
	Page        string `json:"page" validate:"required"`
	Pasaran_id  string `json:"pasaran_id"`
	Prediksi_id int    `json:"prediksi_id"`
}
