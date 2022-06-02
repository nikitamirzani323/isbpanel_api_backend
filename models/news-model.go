package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_newsHome(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_news
	var arraobj []entities.Model_news
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 50
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idnews) as totalnews  "
	sql_selectcount += "FROM " + configs.DB_tbl_trx_news + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(title_news) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(title_news) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "A.idnews , A.title_news, A.descp_news, "
	sql_select += "A.url_news , A.img_news, B.nmcatenews, A.idcatenews, "
	sql_select += "A.createnews, to_char(COALESCE(A.createdatenews,now()), 'YYYY-MM-DD HH24:ii:ss') as creatednews, A.updatenews, to_char(COALESCE(A.updatedatenews,now()), 'YYYY-MM-DD HH24:ii:ss') as updatednews "
	sql_select += "FROM " + configs.DB_tbl_trx_news + " as A "
	sql_select += "JOIN " + configs.DB_tbl_mst_category + " as B ON B.idcatenews = A.idcatenews "
	if search == "" {
		sql_select += "ORDER BY creatednews DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(title_news) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(title_news) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY creatednews DESC  LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idnews_db, idcatenews_db                                              int
			title_news_db, descp_news_db, url_news_db, img_news_db, nmcatenews_db string
			createnews_db, createdatenews_db, updatenews_db, updatedatenews_db    string
		)

		err = row.Scan(
			&idnews_db, &title_news_db, &descp_news_db, &url_news_db, &img_news_db, &nmcatenews_db, &idcatenews_db,
			&createnews_db, &createdatenews_db, &updatenews_db, &updatedatenews_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createnews_db != "" {
			create = createnews_db + ", " + createdatenews_db
		}
		if updatenews_db != "" {
			update = updatenews_db + ", " + updatedatenews_db
		}

		obj.News_id = idnews_db
		obj.News_idcategory = idcatenews_db
		obj.News_category = nmcatenews_db
		obj.News_title = title_news_db
		obj.News_descp = descp_news_db
		obj.News_url = url_news_db
		obj.News_image = img_news_db
		obj.News_create = create
		obj.News_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_news(admin, sdata, title, descp, url, image string, idrecord, category int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_news + ` (
				idnews , idcatenews, title_news, descp_news, url_news, img_news, 
				createnews, createdatenews
			) values (
				$1 ,$2, $3, $4, $5, $6, 
				$7, $8
			)
		`
		field_column := configs.DB_tbl_trx_news + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_news, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			category, title, descp, url, image,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_trx_news + ` 
			SET idcatenews=$1, title_news=$2, descp_news=$3, 
			url_news=$4,img_news=$5,
			updatenews=$6, updatedatenews=$7 
			WHERE idnews=$8 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_news, "UPDATE",
			category, title, descp, url, image,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Delete_news(admin string, idnews int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_news, "idnews", strconv.Itoa(idnews))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_news + ` 
			WHERE idnews=$1 
		`

		flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_trx_news, "DELETE", idnews)

		if flag_delete {
			flag = true
			msg = "Succes"
			log.Println(msg_delete)
		} else {
			log.Println(msg_delete)
		}
	} else {
		msg = "Data Not Found"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Fetch_category() (helpers.Response, error) {
	var obj entities.Model_category
	var arraobj []entities.Model_category
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcatenews , nmcatenews, displaycatenews,statuscategory, 
			createcatenews, COALESCE(createdatecatenews,now()), updatecatenews, COALESCE(updatedatecatenews,now())  
			FROM ` + configs.DB_tbl_mst_category + ` 
			ORDER BY displaycatenews ASC   
		`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcatenews_db, displaycatenews_db                                                  int
			nmcatenews_db, statuscategory_db                                                   string
			createcatenews_db, createdatecatenews_db, updatecatenews_db, updatedatecatenews_db string
		)

		err = row.Scan(
			&idcatenews_db, &nmcatenews_db, &displaycatenews_db, &statuscategory_db,
			&createcatenews_db, &createdatecatenews_db, &updatecatenews_db, &updatedatecatenews_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		statuscss := configs.STATUS_CANCEL
		if createcatenews_db != "" {
			create = createcatenews_db + ", " + createdatecatenews_db
		}
		if updatecatenews_db != "" {
			update = updatecatenews_db + ", " + updatedatecatenews_db
		}
		if statuscategory_db == "Y" {
			statuscss = configs.STATUS_COMPLETE
		}

		obj.Category_id = idcatenews_db
		obj.Category_name = nmcatenews_db
		obj.Category_totalnews = _GetTotalNews(idcatenews_db)
		obj.Category_display = displaycatenews_db
		obj.Category_status = statuscategory_db
		obj.Category_statuscss = statuscss
		obj.Category_create = create
		obj.Category_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_category(admin, name, status, sdata string, idrecord, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_category + ` (
				idcatenews , nmcatenews, displaycatenews, statuscategory, 
				createcatenews, createdatecatenews
			) values (
				$1 ,$2, $3, $4, 
				$5, $6
			)
		`
		field_column := configs.DB_tbl_mst_category + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_category, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			name, display, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_mst_category + ` 
			SET nmcatenews=$1, displaycatenews=$2, statuscategory=$3, 
			updatecatenews=$4, updatedatecatenews=$5 
			WHERE idcatenews=$6 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_category, "UPDATE",
			name, display, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Delete_category(admin string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false
	flag_tblnews := false

	flag = CheckDB(configs.DB_tbl_mst_category, "idcatenews", strconv.Itoa(idrecord))
	if flag {
		flag_tblnews = CheckDB(configs.DB_tbl_trx_news, "idcatenews", strconv.Itoa(idrecord))
		if !flag_tblnews {
			sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_mst_category + ` 
				WHERE idcatenews=$1 
			`
			flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_mst_category, "DELETE", idrecord)

			if flag_delete {
				flag = true
				msg = "Succes"
				log.Println(msg_delete)
			} else {
				log.Println(msg_delete)
			}
		} else {
			msg = "Cannot Delete"
		}
	} else {
		msg = "Data Not Found"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func _GetTotalNews(idrecord int) int {
	con := db.CreateCon()
	ctx := context.Background()
	total := 0

	sql_select := `SELECT
		count(	idnews) as total  
		FROM ` + configs.DB_tbl_trx_news + `  
		WHERE idcatenews = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&total); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return total
}
