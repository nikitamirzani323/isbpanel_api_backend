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

func Fetch_event() (helpers.Response, error) {
	var obj entities.Model_event
	var arraobj []entities.Model_event
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		A.idevent , A.idwebagen, B.nmwebagen, A.nmevent,  A.mindeposit, 
		to_char(COALESCE(A.startevent,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		to_char(COALESCE(A.endevent,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		createevent, to_char(COALESCE(A.createdateevent,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updateevent, to_char(COALESCE(A.updatedateevent,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_event + ` as A 
		JOIN ` + configs.DB_tbl_mst_websiteagen + ` as B ON B.idwebagen = A.idwebagen   
		ORDER BY createdateevent DESC     
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idevent_db, idwebagen_db, mindeposit_db                                int
			nmevent_db, nmwebagen_db, startevent_db, endevent_db                   string
			createevent_db, createdateevent_db, updateevent_db, updatedateevent_db string
		)

		err = row.Scan(&idevent_db, &idwebagen_db,
			&nmwebagen_db, &nmevent_db, &mindeposit_db, &startevent_db, &endevent_db,
			&createevent_db, &createdateevent_db, &updateevent_db, &updatedateevent_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createevent_db != "" {
			create = createevent_db + ", " + createdateevent_db
		}
		if updateevent_db != "" {
			update = updateevent_db + ", " + updatedateevent_db
		}

		obj.Event_id = idevent_db
		obj.Event_idwebagen = idwebagen_db
		obj.Event_nmwebagen = nmwebagen_db
		obj.Event_name = nmevent_db
		obj.Event_startevent = startevent_db
		obj.Event_endevent = endevent_db
		obj.Event_mindeposit = mindeposit_db
		obj.Event_create = create
		obj.Event_update = update
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
func Save_event(
	admin, nmevent, startevent, endevent, sData string,
	idwebagen, mindeposit, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_trx_event + ` (
					idevent , idwebagen, nmevent,  
					startevent , endevent,  mindeposit, 
					createevent, createdateevent
				) values (
					$1, $2, $3, 
					$4, $5,  
					$6, $7
				)
			`
		field_column := configs.DB_tbl_trx_event + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_event, "INSERT",
			tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+strconv.Itoa(idrecord_counter), idwebagen,
			nmevent, startevent, endevent, mindeposit,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_event + `  
				SET idwebagen =$1, nmevent=$2, 
				startevent=$3, endevent=$4, endevent=$5, 
				updateevent=$6, updatedateevent=$7  
				WHERE idevent=$8  
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_event, "UPDATE",
			idwebagen, nmevent, startevent, endevent, mindeposit,
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
func Fetchdetail_event(idevent int) (helpers.Response, error) {
	var obj entities.Model_eventdetail
	var arraobj []entities.Model_eventdetail
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		A.ideventdetail , A.voucher, A.deposit,  
		B.phonemember , B.usernameagen, 
		createeventdetail, to_char(COALESCE(A.createdateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updateeventdetail, to_char(COALESCE(A.updatedateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_event_detail + ` as A 
		JOIN ` + configs.DB_tbl_trx_memberagen + ` as B ON B.idmemberagen = A.idmemberagen    
		WHERE A.idevent=$1 
		ORDER BY A.createdateeventdetail DESC     
	`

	row, err := con.QueryContext(ctx, sql_select, idevent)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			ideventdetail_db, deposit_db                                                                   int
			voucher_db, phonemember_db, usernameagen_db                                                    string
			createeventdetail_db, createdateeventdetail_db, updateeventdetail_db, updatedateeventdetail_db string
		)

		err = row.Scan(&ideventdetail_db, &voucher_db,
			&deposit_db, &phonemember_db, &usernameagen_db,
			&createeventdetail_db, &createdateeventdetail_db, &updateeventdetail_db, &updatedateeventdetail_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createeventdetail_db != "" {
			create = createeventdetail_db + ", " + createdateeventdetail_db
		}
		if updateeventdetail_db != "" {
			update = updateeventdetail_db + ", " + updatedateeventdetail_db
		}

		obj.Eventdetail_iddetail = ideventdetail_db
		obj.Eventdetail_phone = phonemember_db
		obj.Eventdetail_username = usernameagen_db
		obj.Eventdetail_voucher = voucher_db
		obj.Eventdetail_deposit = deposit_db
		obj.Eventdetail_create = create
		obj.Eventdetail_update = update
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
func Savedetail_event(
	admin, sData string,
	idevent, idmemberagen, deposit, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		voucher := "23011717011"
		sql_insert := `
				insert into
				` + configs.DB_tbl_trx_event_detail + ` (
					ideventdetail , idevent, idmemberagen,  
					voucher , deposit,  
					createeventdetail, createdateeventdetail
				) values (
					$1, $2, $3, 
					$4, $5,  
					$6, $7
				)
			`
		field_column := configs.DB_tbl_trx_event_detail + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_event_detail, "INSERT",
			tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+strconv.Itoa(idrecord_counter),
			idevent, idmemberagen, voucher, deposit,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_event_detail + `  
				SET deposit=$1, 
				updateeventdetail=$2, updatedateeventdetail=$3 
				WHERE ideventdetail=$4   
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_event_detail, "UPDATE",
			deposit, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

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
