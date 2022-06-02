package entities

type Model_news struct {
	News_id         int    `json:"news_id"`
	News_idcategory int    `json:"news_idcategory"`
	News_category   string `json:"news_category"`
	News_title      string `json:"news_title"`
	News_descp      string `json:"news_descp"`
	News_url        string `json:"news_url"`
	News_image      string `json:"news_image"`
	News_create     string `json:"news_create"`
	News_update     string `json:"news_update"`
}
type Model_category struct {
	Category_id        int    `json:"category_id"`
	Category_name      string `json:"category_name"`
	Category_totalnews int    `json:"category_totalnews"`
	Category_display   int    `json:"category_display"`
	Category_status    string `json:"category_status"`
	Category_statuscss string `json:"category_statuscss"`
	Category_create    string `json:"category_create"`
	Category_update    string `json:"category_update"`
}

type Controller_news struct {
	News_search string `json:"news_search"`
	News_page   int    `json:"news_page"`
}
type Controller_newssave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	News_idrecord int    `json:"news_id"`
	News_category int    `json:"news_category" validate:"required"`
	News_title    string `json:"news_title" validate:"required"`
	News_descp    string `json:"news_descp" validate:"required"`
	News_url      string `json:"news_url" validate:"required"`
	News_image    string `json:"news_image" validate:"required"`
}
type Controller_newsdelete struct {
	Page    string `json:"page" validate:"required"`
	News_id int    `json:"news_id"`
}
type Controller_categorysave struct {
	Page              string `json:"page" validate:"required"`
	Sdata             string `json:"sdata" validate:"required"`
	Category_idrecord int    `json:"category_id"`
	Category_name     string `json:"category_name" validate:"required"`
	Category_status   string `json:"category_status" `
	Category_display  int    `json:"category_display" validate:"required"`
}
type Controller_categorydelete struct {
	Page              string `json:"page" validate:"required"`
	Category_idrecord int    `json:"category_id" validate:"required"`
}
