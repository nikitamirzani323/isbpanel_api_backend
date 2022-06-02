package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_slideHome() (helpers.Response, error) {
	var obj entities.Model_slide
	var arraobj []entities.Model_slide
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idslide , movieid, url, position,  
			createslide, COALESCE(createdateslide,now())
			FROM ` + configs.DB_tbl_trx_slide + `  
			ORDER BY position ASC  
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idslide_db, movieid_db, position_db int
			url_db                              string
			createslide_db, createdateslide_db  string
		)

		err = row.Scan(&idslide_db, &movieid_db, &url_db, &position_db, &createslide_db, &createdateslide_db)

		helpers.ErrorCheck(err)
		create := ""
		if createslide_db != "" {
			create = createslide_db + ", " + createdateslide_db
		}

		var (
			movietitle_db string
		)
		sql_selectmovie := `SELECT 
			movietitle 
			FROM ` + configs.DB_tbl_trx_movie + ` 
			WHERE movieid = $1 
		`
		row_movie := con.QueryRowContext(ctx, sql_selectmovie, movieid_db)
		switch e_movie := row_movie.Scan(&movietitle_db); e_movie {
		case sql.ErrNoRows:
		case nil:
		default:
			helpers.ErrorCheck(e_movie)
		}

		obj.Slide_id = idslide_db
		obj.Slide_movieid = movieid_db
		obj.Slide_movietitle = movietitle_db
		obj.Slide_urlimage = url_db
		obj.Slide_position = position_db
		obj.Slide_create = create
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
func Save_slide(admin, sdata, url string, idmovie, position int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_slide + ` (
				idslide , movieid, url, position,
				createslide, createdateslide
			) values (
				$1 ,$2, $3, $4,  
				$5, $6
			)
		`
		field_column := configs.DB_tbl_trx_slide + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_slide, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			idmovie, url, position,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
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
func Delete_slider(admin string, idslide int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_slide, "idslide", strconv.Itoa(idslide))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_slide + ` 
			WHERE idslide=$1 
		`
		stmt_delete, e_delete := con.PrepareContext(ctx, sql_delete)
		helpers.ErrorCheck(e_delete)
		defer stmt_delete.Close()
		rec_delete, e_delete := stmt_delete.ExecContext(ctx, idslide)

		helpers.ErrorCheck(e_delete)
		delete, e := rec_delete.RowsAffected()
		helpers.ErrorCheck(e)
		if delete > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di delete")
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
