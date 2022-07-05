package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_tafsirmimpiHome(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_tafsirmimpi
	var arraobj []entities.Model_tafsirmimpi
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
	sql_selectcount += "COUNT(idtafsirmimpi) as total  "
	sql_selectcount += "FROM " + configs.DB_tbl_mst_tafsirmimpi + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(mimpi) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(artimimpi) LIKE '%" + strings.ToLower(search) + "%' "
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
	sql_select += "idtafsirmimpi , mimpi, artimimpi, angka2d, angka3d, angka4d,statustafsirmimpi, "
	sql_select += "createtafsirmimpi , to_char(COALESCE(createdatetafsirmimpi,now()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "updatetafsirmimpi, to_char(COALESCE(updatedatetafsirmimpi,now()), 'YYYY-MM-DD HH24:MI:SS') "
	sql_select += "FROM " + configs.DB_tbl_mst_tafsirmimpi + " "
	if search == "" {
		sql_select += "ORDER BY createtafsirmimpi DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(mimpi) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(artimimpi) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createtafsirmimpi DESC LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtafsirmimpi_db                                                                               int
			mimpi_db, artimimpi_db, angka2d_db, angka3d_db, angka4d_db, statustafsirmimpi_db               string
			createtafsirmimpi_db, createdatetafsirmimpi_db, updatetafsirmimpi_db, updatedatetafsirmimpi_db string
		)

		err = row.Scan(
			&idtafsirmimpi_db, &mimpi_db, &artimimpi_db, &angka2d_db, &angka3d_db, &angka4d_db, &statustafsirmimpi_db,
			&createtafsirmimpi_db, &createdatetafsirmimpi_db, &updatetafsirmimpi_db, &updatedatetafsirmimpi_db)

		helpers.ErrorCheck(err)
		statuscss := configs.STATUS_CANCEL
		create := ""
		update := ""
		if createtafsirmimpi_db != "" {
			create = createtafsirmimpi_db + ", " + createdatetafsirmimpi_db
		}
		if updatetafsirmimpi_db != "" {
			update = updatetafsirmimpi_db + ", " + updatedatetafsirmimpi_db
		}
		if statustafsirmimpi_db == "Y" {
			statuscss = configs.STATUS_RUNNING
		}

		obj.Tafsirmimpi_id = idtafsirmimpi_db
		obj.Tafsirmimpi_mimpi = mimpi_db
		obj.Tafsirmimpi_artimimpi = artimimpi_db
		obj.Tafsirmimpi_angka2d = angka2d_db
		obj.Tafsirmimpi_angka3d = angka3d_db
		obj.Tafsirmimpi_angka4d = angka4d_db
		obj.Tafsirmimpi_status = statustafsirmimpi_db
		obj.Tafsirmimpi_statuscss = statuscss
		obj.Tafsirmimpi_create = create
		obj.Tafsirmimpi_update = update
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
func Save_tafsirmimpi(admin, mimpi, artimimpi, angka2d, angka3d, angka4d, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_tafsirmimpi + ` (
				idtafsirmimpi , mimpi, artimimpi, angka2d, angka3d, angka4d, statustafsirmimpi, 
				createtafsirmimpi, createdatetafsirmimpi
			) values (
				$1 ,$2, $3, $4, $5, $6, $7,
				$8, $9
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_mst_tafsirmimpi + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			mimpi, artimimpi, angka2d, angka3d, angka4d, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))
		helpers.ErrorCheck(e_newrecord)
		insert, e := res_newrecord.RowsAffected()
		helpers.ErrorCheck(e)
		if insert > 0 {
			msg = "Succes"
			log.Println("Data Berhasil di save")
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_tafsirmimpi + `  
				SET mimpi=$1,artimimpi=$2, angka2d=$3, angka3d=$4, angka4d=$5, statustafsirmimpi=$6,
				updatetafsirmimpi=$7, updatedatetafsirmimpi=$8 
				WHERE idtafsirmimpi =$9 
			`
		stmt_record, e := con.PrepareContext(ctx, sql_update)
		helpers.ErrorCheck(e)
		rec_record, e_record := stmt_record.ExecContext(
			ctx,
			mimpi, artimimpi, angka2d, angka3d, angka4d, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			idrecord)
		helpers.ErrorCheck(e_record)
		update_record, e_record := rec_record.RowsAffected()
		helpers.ErrorCheck(e_record)

		defer stmt_record.Close()
		if update_record > 0 {
			msg = "Succes"
			log.Printf("Update PASARAN Success : %d\n", idrecord)
		} else {
			log.Println("Update PASARAN failed")
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
