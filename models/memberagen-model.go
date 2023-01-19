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
	"github.com/buger/jsonparser"
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

		//WEBSITE AGEN
		var objwebsiteagen entities.Model_memberagen
		var arraobjwebsiteagen []entities.Model_memberagen
		sql_selectwebsiteagen := `SELECT 
			A.idwebagen, A.usernameagen, B.nmwebagen  
			FROM ` + configs.DB_tbl_trx_memberagen + ` as A 
			JOIN ` + configs.DB_tbl_mst_websiteagen + ` as B ON B.idwebagen = A.idwebagen 
			WHERE A.phonemember = $1   
		`
		row_websiteagen, err_websiteagen := con.QueryContext(ctx, sql_selectwebsiteagen, phonemember_db)
		helpers.ErrorCheck(err_websiteagen)
		for row_websiteagen.Next() {
			var (
				idwebagen_db                  int
				usernameagen_db, nmwebagen_db string
			)
			err_websiteagen = row_websiteagen.Scan(&idwebagen_db, &usernameagen_db, &nmwebagen_db)
			objwebsiteagen.Memberagen_idwebagen = idwebagen_db
			objwebsiteagen.Memberagen_username = usernameagen_db
			objwebsiteagen.Memberagen_website = nmwebagen_db
			arraobjwebsiteagen = append(arraobjwebsiteagen, objwebsiteagen)
		}

		obj.Member_phone = phonemember_db
		obj.Member_name = nmmember_db
		obj.Member_agen = arraobjwebsiteagen
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
func Fetch_memberSelect(idwebagen int) (helpers.Response, error) {
	var obj entities.Model_memberagenselect
	var arraobj []entities.Model_memberagenselect
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		A.idmemberagen, A.phonemember, A.usernameagen, C.nmmember, 
		B.nmwebagen 
		FROM ` + configs.DB_tbl_trx_memberagen + ` as A   
		JOIN ` + configs.DB_tbl_mst_websiteagen + ` as B ON B.idwebagen = A.idwebagen    
		JOIN ` + configs.DB_tbl_trx_member + ` as C ON C.phonemember = A.phonemember    
		WHERE A.idwebagen = $1 
		ORDER BY C.nmmember ASC      
	`

	row, err := con.QueryContext(ctx, sql_select, idwebagen)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmemberagen_db                                            int
			phonemember_db, usernameagen_db, nmmember_db, nmwebagen_db string
		)

		err = row.Scan(&idmemberagen_db, &phonemember_db,
			&usernameagen_db, &nmmember_db, &nmwebagen_db)

		helpers.ErrorCheck(err)

		obj.Memberagen_id = idmemberagen_db
		obj.Memberagen_username = usernameagen_db
		obj.Memberagen_phone = phonemember_db
		obj.Memberagen_website = nmwebagen_db
		obj.Memberagen_name = nmmember_db
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
	admin, phone, nama, listagen, sData string,
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

				//WEBSITE AGEN
				_save_memberagen(admin, phone, listagen)
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
			_delete_memberagen(phone)
			_save_memberagen(admin, phone, listagen)
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
func _save_memberagen(admin, phone, listagen string) {
	tglnow, _ := goment.New()
	jsonsource := []byte(listagen)
	jsonparser.ArrayEach(jsonsource, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agen_idwebsite, _ := jsonparser.GetInt(value, "agen_idwebsite")
		agen_username, _ := jsonparser.GetString(value, "agen_username")

		sql_insertagen := `
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
		flag_insertagen, msg_insertagen := Exec_SQL(sql_insertagen, configs.DB_tbl_trx_event, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), agen_idwebsite,
			phone, agen_username,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if !flag_insertagen {
			log.Println(msg_insertagen)
		}
	})
}
func _delete_memberagen(phone string) {
	sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_trx_memberagen + ` 
				WHERE phonemember=$1  
			`
	flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_trx_memberagen, "DELETE", phone)

	if !flag_delete {
		log.Println(msg_delete)
	}
}
func _GetMemberAgen(idrecord int) (string, string) {
	con := db.CreateCon()
	ctx := context.Background()
	phonemember := ""
	usernameagen := ""

	sql_select := `SELECT
		phonemember, usernameagen   
		FROM ` + configs.DB_tbl_trx_memberagen + `  
		WHERE idmemberagen = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&phonemember, &usernameagen); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return phonemember, usernameagen
}
