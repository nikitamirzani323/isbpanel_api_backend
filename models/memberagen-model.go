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

func Fetch_member() (helpers.Response, error) {
	var obj entities.Model_member
	var arraobj []entities.Model_member
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		phonemember, nmmember,   
		createmember, to_char(COALESCE(createdatemember,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updatemember, to_char(COALESCE(updatedatemember,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_member + `  
		ORDER BY createdatemember DESC     
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			phonemember_db, nmmember_db                                                string
			createmember_db, createdatemember_db, updatemember_db, updatedatemember_db string
		)

		err = row.Scan(&phonemember_db, &nmmember_db,
			&createmember_db, &createdatemember_db, &updatemember_db, &updatedatemember_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createmember_db != "" {
			create = createmember_db + ", " + createdatemember_db
		}
		if updatemember_db != "" {
			update = updatemember_db + ", " + updatedatemember_db
		}

		obj.Member_phone = phonemember_db
		obj.Member_name = nmmember_db
		obj.Member_create = create
		obj.Member_update = update
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
func Fetch_memberagen(phone string) (helpers.Response, error) {
	var obj entities.Model_memberagen
	var arraobj []entities.Model_memberagen
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idmemberagen, idwebagen, usernameagen,   
		creatememberagen, to_char(COALESCE(createdatememberagen,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updatememberagen, to_char(COALESCE(updatedatememberagen,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_memberagen + `  
		WHERE phonemember=$1 
		ORDER BY createdatememberagen DESC     
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmemberagen_db, idwebagen_db                                                              int
			usernameagen_db                                                                            string
			creatememberagen_db, createdatememberagen_db, updatememberagen_db, updatedatememberagen_db string
		)

		err = row.Scan(&idmemberagen_db, &idwebagen_db,
			&usernameagen_db,
			&creatememberagen_db, &createdatememberagen_db, &updatememberagen_db, &updatedatememberagen_db, &phone)

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
		obj.Memberagen_username = usernameagen_db
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
func Save_member(
	admin, phone, nama, sData string,
	idrecord string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_trx_member, "phone", phone)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_trx_member + ` (
					phonemember, nmmember,  
					createmember, createdatemember
				) values (
					$1, $2, 
					$3, $4
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_member, "INSERT",
				phone, nama,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}

	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_member + `  
				SET nmmember=$1, 
				updatemember=$2, updatedatemember=$3 
				WHERE phonemember =$4 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_member, "UPDATE",
			nama, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
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
func Save_memberagen(
	admin, username, phone, sData string,
	idwebagen, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_trx_memberagen + ` (
					idmemberagen, idwebagen, phonemember, usernameagen, 
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
			phone, username,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			log.Println(msg_insert)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Delete_memberagen(phone string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_trx_memberagen, "idmemberagen", strconv.Itoa(idrecord), "phonemember", phone)
	if flag {
		sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_trx_memberagen + ` 
				WHERE idmemberagen=$1 AND phonemember=$2 
			`
		flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_trx_memberagen, "DELETE", idrecord, phone)

		if flag_delete {
			flag = true
			msg = "Succes"
		} else {
			log.Println(msg_delete)
		}
	} else {
		msg = "Data Not Found"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
