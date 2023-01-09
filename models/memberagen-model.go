package models

import (
	"context"
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

func Fetch_memberagen() (helpers.Response, error) {
	var obj entities.Model_memberagen
	var arraobj []entities.Model_memberagen
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idmemberagen, usernamememberagen , idwebagen, nmmember,  
		creatememberagen, to_char(COALESCE(createdatememberagen,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updatememberagen, to_char(COALESCE(updatedatememberagen,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_memberagen + `  
		ORDER BY createdatememberagen DESC     
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmemberagen_db, idwebagen_db                                                              int
			usernamememberagen_db, nmmember_db                                                         string
			creatememberagen_db, createdatememberagen_db, updatememberagen_db, updatedatememberagen_db string
		)

		err = row.Scan(&idmemberagen_db, &idwebagen_db,
			&usernamememberagen_db, &nmmember_db,
			&creatememberagen_db, &createdatememberagen_db, &updatememberagen_db, &updatedatememberagen_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if creatememberagen_db != "" {
			create = creatememberagen_db + ", " + createdatememberagen_db
		}
		if updatememberagen_db != "" {
			update = updatememberagen_db + ", " + updatedatememberagen_db
		}

		obj.Memberagen_id = idmemberagen_db
		obj.Memberagen_idwebagen = idwebagen_db
		obj.Memberagen_username = usernamememberagen_db
		obj.Memberagen_name = nmmember_db
		obj.Memberagen_create = create
		obj.Memberagen_update = update
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
func Save_memberagen(
	admin, username, nama, sData string,
	idwebagen, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_trx_memberagen + ` (
					idmemberagen, usernamememberagen , idwebagen, nmmember, 
					creatememberagen, createdatememberagen
				) values (
					$1, $2, $3, $4,
					$5, $6
				)
			`
		field_column := configs.DB_tbl_trx_memberagen + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_event, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idwebagen,
			username, idwebagen, nama,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_memberagen + `  
				SET idwebagen =$1, nmmember=$2, 
				updatememberagen=$3, updatedatememberagen=$4 
				WHERE idmemberagen=$5   
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_memberagen, "UPDATE",
			idwebagen, nama,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
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
