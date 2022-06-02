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

func Fetch_websiteagenHome(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_websiteagen
	var arraobj []entities.Model_websiteagen
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
	sql_selectcount += "COUNT(idwebagen) as totalwebagen  "
	sql_selectcount += "FROM " + configs.DB_tbl_mst_websiteagen + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(nmwebagen) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nmwebagen) LIKE '%" + strings.ToLower(search) + "%' "
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
	sql_select += "idwebagen ,nmwebagen, statuswebagen, "
	sql_select += "createwebagen, to_char(COALESCE(createdatewebagen,now()), 'YYYY-MM-DD HH24:ii:ss') as createdatewebagen, updatewebagen, to_char(COALESCE(updatedatewebagen,now()), 'YYYY-MM-DD HH24:ii:ss') as updatedatewebagen "
	sql_select += "FROM " + configs.DB_tbl_mst_websiteagen + " "
	if search == "" {
		sql_select += "ORDER BY createwebagen DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(nmwebagen) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(nmwebagen) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createwebagen DESC  LIMIT " + strconv.Itoa(perpage)
	}
	log.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idwebagen_db                                                                   int
			nmwebagen_db, statuswebagen_db                                                 string
			createwebagen_db, createdatewebagen_db, updatewebagen_db, updatedatewebagen_db string
		)

		err = row.Scan(
			&idwebagen_db, &nmwebagen_db, &statuswebagen_db,
			&createwebagen_db, &createdatewebagen_db, &updatewebagen_db, &updatedatewebagen_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createwebagen_db != "" {
			create = createwebagen_db + ", " + createdatewebagen_db
		}
		if updatewebagen_db != "" {
			update = updatewebagen_db + ", " + updatedatewebagen_db
		}

		obj.Websiteagen_id = idwebagen_db
		obj.Websiteagen_name = nmwebagen_db
		obj.Websiteagen_status = statuswebagen_db
		obj.Websiteagen_create = create
		obj.Websiteagen_update = update
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
func Save_websiteagen(admin, sdata, name, status string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_websiteagen + ` (
				idwebagen , nmwebagen, statuswebagen,  
				createwebagen, createdatewebagen  
			) values (
				$1 ,$2, $3, 
				$4, $5  
			)
		`
		field_column := configs.DB_tbl_mst_websiteagen + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_websiteagen, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_mst_websiteagen + ` 
			SET nmwebagen=$1, statuswebagen=$2,  
			updatewebagen=$3, updatedatewebagen=$4 
			WHERE idwebagen=$5  
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_websiteagen, "UPDATE",
			name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
