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

func Fetch_event() (helpers.Response, error) {
	var obj entities.Model_event
	var arraobj []entities.Model_event
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		A.idevent , A.idwebagen, B.nmwebagen, A.nmevent,  A.mindeposit, A.money_in, A.money_out, 
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
			money_in_db, money_out_db                                              float32
			nmevent_db, nmwebagen_db, startevent_db, endevent_db                   string
			createevent_db, createdateevent_db, updateevent_db, updatedateevent_db string
		)

		err = row.Scan(&idevent_db, &idwebagen_db,
			&nmwebagen_db, &nmevent_db, &mindeposit_db, &money_in_db, &money_out_db, &startevent_db, &endevent_db,
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
		obj.Event_money_in = int(money_in_db)
		obj.Event_money_out = int(money_out_db)
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
					money_in, money_out, 
					createevent, createdateevent
				) values (
					$1, $2, $3, 
					$4, $5, $6, 
					$7, $8,
					$9, $10 
				)
			`
		field_column := configs.DB_tbl_trx_event + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_event, "INSERT",
			tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+strconv.Itoa(idrecord_counter), idwebagen,
			nmevent, startevent, endevent, mindeposit, 0, 0,
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
func Fetchdetail_event(idevent, idmemberagen int) (helpers.Response, error) {
	var obj entities.Model_eventdetail
	var arraobj []entities.Model_eventdetail
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "A.ideventdetail , A.voucher, A.deposit, A.statuseventdetail, "
	sql_select += "B.phonemember , B.usernameagen, "
	sql_select += "createeventdetail, to_char(COALESCE(A.createdateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS'),  "
	sql_select += "updateeventdetail, to_char(COALESCE(A.updatedateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS')   "
	sql_select += "FROM " + configs.DB_tbl_trx_event_detail + "  as A "
	sql_select += "JOIN " + configs.DB_tbl_trx_memberagen + "  as B ON B.idmemberagen = A.idmemberagen "
	sql_select += "WHERE A.idevent=$1 "
	if idmemberagen > 0 {
		sql_select += "AND A.idmemberagen='" + strconv.Itoa(idmemberagen) + "' "
	}
	sql_select += "ORDER BY A.createdateeventdetail DESC "
	log.Println(sql_select)

	row, err := con.QueryContext(ctx, sql_select, idevent)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			ideventdetail_db, deposit_db                                                                   int
			voucher_db, phonemember_db, usernameagen_db, statuseventdetail_db                              string
			createeventdetail_db, createdateeventdetail_db, updateeventdetail_db, updatedateeventdetail_db string
		)

		err = row.Scan(&ideventdetail_db, &voucher_db,
			&deposit_db, &statuseventdetail_db, &phonemember_db, &usernameagen_db,
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
		obj.Eventdetail_status = statuseventdetail_db
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
func Fetchdetailwinner_event(idevent int) (helpers.Response, error) {
	var obj entities.Model_eventdetail
	var arraobj []entities.Model_eventdetail
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "A.ideventdetail , A.voucher, A.deposit, A.statuseventdetail, "
	sql_select += "B.phonemember , B.usernameagen, "
	sql_select += "createeventdetail, to_char(COALESCE(A.createdateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS'),  "
	sql_select += "updateeventdetail, to_char(COALESCE(A.updatedateeventdetail,now()), 'YYYY-MM-DD HH24:MI:SS')   "
	sql_select += "FROM " + configs.DB_tbl_trx_event_detail + "  as A "
	sql_select += "JOIN " + configs.DB_tbl_trx_memberagen + "  as B ON B.idmemberagen = A.idmemberagen "
	sql_select += "WHERE A.idevent=$1 "
	sql_select += "AND A.statuseventdetail!='' "
	sql_select += "ORDER BY A.createdateeventdetail DESC "

	row, err := con.QueryContext(ctx, sql_select, idevent)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			ideventdetail_db, deposit_db                                                                   int
			voucher_db, phonemember_db, usernameagen_db, statuseventdetail_db                              string
			createeventdetail_db, createdateeventdetail_db, updateeventdetail_db, updatedateeventdetail_db string
		)

		err = row.Scan(&ideventdetail_db, &voucher_db,
			&deposit_db, &statuseventdetail_db, &phonemember_db, &usernameagen_db,
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
		obj.Eventdetail_status = statuseventdetail_db
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
func Fetchdetailgroup_event(idevent int) (helpers.Response, error) {
	var obj entities.Model_eventdetailgroup
	var arraobj []entities.Model_eventdetailgroup
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idmemberagen, SUM(deposit) as totaldeposit, COUNT(deposit) as totalkupon    
		FROM ` + configs.DB_tbl_trx_event_detail + ` 
		WHERE idevent=$1 
		GROUP BY idmemberagen 
		ORDER BY totaldeposit DESC     
	`

	row, err := con.QueryContext(ctx, sql_select, idevent)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idmemberagen_db                int
			totalkupon_db, totaldeposit_db float32
		)

		err = row.Scan(&idmemberagen_db, &totaldeposit_db, &totalkupon_db)

		helpers.ErrorCheck(err)
		phone, username := _GetMemberAgen(idmemberagen_db)
		obj.Eventdetailgroup_idmember = idmemberagen_db
		obj.Eventdetailgroup_username = username
		obj.Eventdetailgroup_phone = phone
		obj.Eventdetailgroup_deposit = int(totaldeposit_db)
		obj.Eventdetailgroup_voucher = int(totalkupon_db)
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
	idevent, idmemberagen, qty, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		mindeposit := _GetEvent(idevent)
		for i := 0; i < qty; i++ {
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
			voucher := strconv.Itoa(idevent) + tglnow.Format("MM") + tglnow.Format("DD") + tglnow.Format("HH") + strconv.Itoa(idrecord_counter)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_event_detail, "INSERT",
				tglnow.Format("YY")+tglnow.Format("MM")+tglnow.Format("DD")+strconv.Itoa(idrecord_counter),
				idevent, idmemberagen, voucher, mindeposit,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				log.Println(msg_insert)
			}
		}
		_updateEvent(admin, idevent)
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Updatedetailstatus_event(
	admin, status, sData string,
	idevent, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_trx_event_detail, "ideventdetail", strconv.Itoa(idrecord), "statuseventdetail", "")

	if flag {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_event_detail + `  
				SET statuseventdetail =$1, 
				updateeventdetail=$2, updatedateeventdetail=$3  
				WHERE ideventdetail=$4  
				AND idevent=$5
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_event_detail, "UPDATE",
			status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord, idevent)

		if !flag_update {
			log.Println(msg_update)
		} else {
			msg = "Succes"
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _updateEvent(admin string, idrecord int) {
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	money_in := 0
	sql_select := `SELECT 
		SUM(deposit) as totaldeposit 
		FROM ` + configs.DB_tbl_trx_event_detail + ` 
		WHERE idevent=$1 
		GROUP BY idevent   
	`

	row, err := con.QueryContext(ctx, sql_select, idrecord)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			totaldeposit_db float32
		)

		err = row.Scan(&totaldeposit_db)
		money_in = int(totaldeposit_db)
		helpers.ErrorCheck(err)
	}
	defer row.Close()
	if money_in > 0 {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_event + `  
				SET money_in =$1, 
				updateevent=$2, updatedateevent=$3  
				WHERE idevent=$4  
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_event, "UPDATE",
			money_in,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if !flag_update {
			log.Println(msg_update)
		}
	}
}
func _GetEvent(idevent int) int {
	con := db.CreateCon()
	ctx := context.Background()

	var (
		mindeposit_db int
	)

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "mindeposit "
	sql_select += "FROM " + configs.DB_tbl_trx_event + " "
	sql_select += "WHERE idevent = $1 "

	log.Println(sql_select)
	row := con.QueryRowContext(ctx, sql_select, idevent)
	switch e := row.Scan(&mindeposit_db); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return mindeposit_db
}
