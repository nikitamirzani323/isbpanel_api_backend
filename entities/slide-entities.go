package entities

type Model_slide struct {
	Slide_id         int    `json:"slide_id"`
	Slide_movieid    int    `json:"slide_movieid"`
	Slide_movietitle string `json:"slide_movietitle"`
	Slide_urlimage   string `json:"slide_urlimage"`
	Slide_position   int    `json:"slide_position"`
	Slide_create     string `json:"slide_create"`
}

type Controller_slidesave struct {
	Page           string `json:"page" validate:"required"`
	Sdata          string `json:"sdata" validate:"required"`
	Slide_movieid  int    `json:"slide_movieid" validate:"required"`
	Slide_urlimage string `json:"slide_urlimage" validate:"required"`
	Slide_position int    `json:"slide_position" validate:"required"`
}
type Controller_slidedelete struct {
	Page     string `json:"page" validate:"required"`
	Slide_id int    `json:"slide_id"`
}
