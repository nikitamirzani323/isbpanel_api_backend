package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_backend/helpers"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_albumHome(page int) (helpers.Responsemovie, error) {
	var obj entities.Model_album
	var arraobj []entities.Model_album
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
	sql_selectcount += "COUNT(idalbum) as totalalbum  "
	sql_selectcount += "FROM " + configs.DB_tbl_trx_albumcloudflare + "  "

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "idalbum , idcloudflare, filenamecloudflare, signedurl, link1, COALESCE(createdatealbum,now()) "
	sql_select += "FROM " + configs.DB_tbl_trx_albumcloudflare + " "
	sql_select += "ORDER BY idalbum DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idalbum_db                                                     int
			idcloudflare_db, filenamecloudflare_db, signedurl_db, link1_db string
			createdatealbum_db                                             string
		)

		err = row.Scan(&idalbum_db, &idcloudflare_db, &filenamecloudflare_db,
			&signedurl_db, &link1_db, &createdatealbum_db)

		helpers.ErrorCheck(err)
		create := ""
		if createdatealbum_db != "" {
			create = createdatealbum_db
		}
		movieid, movietitle, enabled_db := _GetMovie(link1_db)
		status := "HIDE"
		statuscss := configs.STATUS_CANCEL
		if enabled_db > 0 {
			status = "SHOW"
			statuscss = configs.STATUS_RUNNING
		}
		obj.Album_id = idalbum_db
		obj.Album_idcloud = idcloudflare_db
		obj.Album_name = filenamecloudflare_db
		obj.Album_signed = signedurl_db
		obj.Album_varian = link1_db
		obj.Album_movieid = movieid
		obj.Album_movie = movietitle
		obj.Album_moviestatus = status
		obj.Album_moviestatuscss = statuscss
		obj.Album_createdate = create
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
func Save_album(admin, datacloudflare, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		json := []byte(datacloudflare)
		jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			album_id, _ := jsonparser.GetString(value, "album_id")
			album_filename, _ := jsonparser.GetString(value, "album_filename")
			album_signed, _ := jsonparser.GetString(value, "album_signed")
			album_link_0, _ := jsonparser.GetString(value, "album_link_0")

			flag_check := CheckDB(configs.DB_tbl_trx_albumcloudflare, "idcloudflare", album_id)
			log.Printf("%s - %s - %s - %s - %t", album_id, album_filename, album_signed, album_link_0, flag_check)
			if !flag_check {
				sql_insert := `
					insert into
					` + configs.DB_tbl_trx_albumcloudflare + ` (
						idalbum , idcloudflare, filenamecloudflare, signedurl, link1,createdatealbum
					) values (
						$1, $2, $3, $4, $5, $6
					)
				`
				field_column := configs.DB_tbl_trx_albumcloudflare + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_albumcloudflare, "INSERT",
					tglnow.Format("YY")+strconv.Itoa(idrecord_counter), album_id, album_filename,
					album_signed, album_link_0,
					tglnow.Format("YYYY-MM-DD HH:mm:ss"))

				if flag_insert {
					flag = true
					msg = "Succes"
					log.Println(msg_insert)
				} else {
					log.Println(msg_insert)
				}
			}
		})
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
func Update_album(admin, signedurl string, idrecord int) bool {
	flag := false

	sql_update := `
			UPDATE 
			` + configs.DB_tbl_trx_albumcloudflare + `  
			SET signedurl =$1 
			WHERE idalbum =$2 
		`

	flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_albumcloudflare, "UPDATE", signedurl, idrecord)

	if flag_update {
		flag = true
		log.Println(msg_update)
	} else {
		log.Println(msg_update)
	}

	return flag
}
func Delete_album(admin string, idrecord, idmovie int) bool {
	tglnow, _ := goment.New()
	flag := false

	sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_trx_albumcloudflare + ` 
				WHERE idalbum=$1 
			`
	flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_trx_albumcloudflare, "DELETE", idrecord)

	if flag_delete {
		flag = true
		log.Println(msg_delete)
		flag_movie := CheckDB(configs.DB_tbl_trx_movie, "movieid", strconv.Itoa(idmovie))

		if flag_movie {
			sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_movie + `  
				SET urlthumbnail =$1 , updatemovie=$2, updatedatemovie=$3  
				WHERE movieid =$4 
			`

			flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_movie, "UPDATE",
				"https://imagedelivery.net/W-Usm3AjeE17sxpltvGRNA/fd0287a2-353d-4b47-9a6c-9c8df2ab3f00/public",
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmovie)
			if flag_update {
				flag = true
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		}
	} else {
		log.Println(msg_delete)
	}

	return flag
}
func _GetMovie(urlthum string) (int, string, int) {
	con := db.CreateCon()
	ctx := context.Background()
	movietitle := ""
	movieid := 0
	enabled := 0

	sql_select := `SELECT
		movieid, movietitle, enabled   
		FROM ` + configs.DB_tbl_trx_movie + `  
		WHERE urlthumbnail = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, urlthum)
	switch e := row.Scan(&movieid, &movietitle, &enabled); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return movieid, movietitle, enabled
}
