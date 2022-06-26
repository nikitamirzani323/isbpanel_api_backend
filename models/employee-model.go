package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_employeeHome() (helpers.ResponseEmployee, error) {
	var obj entities.Model_employee
	var arraobj []entities.Model_employee
	var objdepart entities.Model_listdepartement
	var arraobjdepart []entities.Model_listdepartement
	var res helpers.ResponseEmployee
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			A.username , A.iddepartement, B.nmdepartement,   
			A.nmemployee , A.phoneemployee, A.statusemployee,  
			A.createemployee, to_char(COALESCE(A.createdateemployee,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			A.updateemployee, to_char(COALESCE(A.updatedateemployee,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_employee + ` as A 
			JOIN ` + configs.DB_tbl_mst_departement + ` as B ON B.iddepartement = A.iddepartement  
			ORDER BY A.createdateemployee DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_db, iddepartement_db, nmdepartement_db                                    string
			nmemployee_db, phoneemployee_db, statusemployee_db                                 string
			createemployee_db, createdateemployee_db, updateemployee_db, updatedateemployee_db string
		)

		err = row.Scan(&username_db, &iddepartement_db, &nmdepartement_db,
			&nmemployee_db, &phoneemployee_db, &statusemployee_db,
			&createemployee_db, &createdateemployee_db, &updateemployee_db, &updatedateemployee_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createemployee_db != "" {
			create = createemployee_db + ", " + createdateemployee_db
		}
		if updateemployee_db != "" {
			update = updateemployee_db + ", " + updatedateemployee_db
		}

		obj.Employee_username = username_db
		obj.Employee_iddepart = iddepartement_db
		obj.Employee_nmdepart = nmdepartement_db
		obj.Employee_name = nmemployee_db
		obj.Employee_phone = phoneemployee_db
		obj.Employee_status = statusemployee_db
		obj.Employee_create = create
		obj.Employee_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	sql_selectdepart := `SELECT 
			iddepartement, nmdepartement 
			FROM ` + configs.DB_tbl_mst_departement + ` 
			ORDER BY nmdepartement ASC    
	`
	rowdepart, errdepart := con.QueryContext(ctx, sql_selectdepart)
	helpers.ErrorCheck(errdepart)
	for rowdepart.Next() {
		var (
			iddepartement_db, nmdepartement_db string
		)

		errdepart = rowdepart.Scan(&iddepartement_db, &nmdepartement_db)

		helpers.ErrorCheck(errdepart)

		objdepart.Departement_id = iddepartement_db
		objdepart.Departement_name = nmdepartement_db
		arraobjdepart = append(arraobjdepart, objdepart)
		msg = "Success"
	}
	defer row.Close()
	defer rowdepart.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Listdepartement = arraobjdepart
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Fetch_employeeByDepartement(iddepart string) (helpers.Response, error) {
	var obj entities.Model_employeebydepart
	var arraobj []entities.Model_employeebydepart
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			username , nmemployee 
			FROM ` + configs.DB_tbl_mst_employee + `  
			WHERE iddepartement=$1 
			AND statusemployee='Y' 
			ORDER BY nmemployee ASC    
	`

	row, err := con.QueryContext(ctx, sql_select, iddepart)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			username_db, nmemployee_db string
		)

		err = row.Scan(&username_db, &nmemployee_db)

		helpers.ErrorCheck(err)

		obj.Employee_username = username_db
		obj.Employee_name = nmemployee_db
		obj.Employee_deposit = _GetSalesStatus(username_db, "DEPOSIT", "")
		obj.Employee_reject = _GetSalesStatus(username_db, "REJECT", "")
		obj.Employee_noanswer = _GetSalesStatus(username_db, "NOANSWER", "")
		obj.Employee_invalid = _GetSalesStatus(username_db, "INVALID", "")
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
func Fetch_employeeBySalesPerformance(iddepart, username, startdate, enddate string) (helpers.Response, error) {
	var obj entities.Model_employeebysalesperform
	var arraobj []entities.Model_employeebysalesperform
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()
	tglnow_start, _ := goment.New(startdate)
	tglnow_end, _ := goment.New(enddate)

	//LIST DEPOSIT
	var obj_listdeposit entities.Model_crmmemberlistdeposit
	var arraobj_listdeposit []entities.Model_crmmemberlistdeposit
	sql_select := ""
	sql_select += "SELECT "
	sql_select += "B.phone , B.nama, B.source,  "
	sql_select += "C.nmwebagen, A.iduseragen, A.deposit,   "
	sql_select += "A.updatecrmsales, to_char(COALESCE(A.updatedatecrmsales,now()), 'YYYY-MM-DD HH24:MI:SS')   "
	sql_select += "FROM " + configs.DB_tbl_trx_crmsales + " as A "
	sql_select += "JOIN " + configs.DB_tbl_trx_usersales + " as B ON B.phone = A.phone "
	sql_select += "JOIN " + configs.DB_tbl_mst_websiteagen + " as C ON C.idwebagen = A.idwebagen "
	sql_select += "WHERE A.username=$1  "
	sql_select += "AND A.statuscrmsales_satu='VALID'  "
	sql_select += "AND A.statuscrmsales_dua='DEPOSIT'  "
	if startdate != "" {
		sql_select += "AND A.updatedatecrmsales >='" + tglnow_start.Format("YYYY-MM-DD") + " 00:00:00' "
		sql_select += "AND A.updatedatecrmsales <='" + tglnow_end.Format("YYYY-MM-DD") + " 23:59:59' "
		sql_select += "ORDER BY A.updatedatecrmsales DESC   "
	} else {
		sql_select += "ORDER BY A.updatedatecrmsales DESC LIMIT 70  "
	}

	row, err := con.QueryContext(ctx, sql_select, username)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			phone_db, nama_db, source_db             string
			nmwebagen_db, iduseragen_db              string
			updatecrmsales_db, updatedatecrmsales_db string
			deposit_db                               float32
		)

		err = row.Scan(&phone_db, &nama_db, &source_db, &nmwebagen_db, &iduseragen_db, &deposit_db,
			&updatecrmsales_db, &updatedatecrmsales_db)

		helpers.ErrorCheck(err)

		obj_listdeposit.Crmsdeposit_phone = phone_db
		obj_listdeposit.Crmsdeposit_nama = nama_db
		obj_listdeposit.Crmsdeposit_source = source_db
		obj_listdeposit.Crmsdeposit_nmwebagen = nmwebagen_db
		obj_listdeposit.Crmsdeposit_iduseragen = iduseragen_db
		obj_listdeposit.Crmsdeposit_deposit = deposit_db
		obj_listdeposit.Crmsdeposit_update = updatecrmsales_db + ", " + updatedatecrmsales_db
		arraobj_listdeposit = append(arraobj_listdeposit, obj_listdeposit)
		msg = "Success"
	}
	defer row.Close()

	//LIST NOANSWER
	var obj_listnoanswer entities.Model_crmmemberlistnoanswer
	var arraobj_listnoanswer []entities.Model_crmmemberlistnoanswer
	sql_select_noanswer := ""
	sql_select_noanswer += "SELECT "
	sql_select_noanswer += "B.phone , B.nama, B.source, A.statuscrmsales_dua, A.notecrmsales,  "
	sql_select_noanswer += "A.updatecrmsales, to_char(COALESCE(A.updatedatecrmsales,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select_noanswer += "FROM " + configs.DB_tbl_trx_crmsales + " as A "
	sql_select_noanswer += "JOIN " + configs.DB_tbl_trx_usersales + " as B ON B.phone = A.phone "
	sql_select_noanswer += "WHERE A.username=$1 "
	sql_select_noanswer += "AND A.statuscrmsales_satu='VALID' "
	sql_select_noanswer += "AND A.statuscrmsales_dua!='DEPOSIT' "
	if startdate != "" {
		sql_select_noanswer += "AND A.updatedatecrmsales >='" + tglnow_start.Format("YYYY-MM-DD") + " 00:00:00' "
		sql_select_noanswer += "AND A.updatedatecrmsales <='" + tglnow_end.Format("YYYY-MM-DD") + " 23:59:59' "
		sql_select_noanswer += "ORDER BY A.updatedatecrmsales DESC   "
	} else {
		sql_select_noanswer += "ORDER BY A.updatedatecrmsales DESC LIMIT 70  "
	}

	row_noanswer, err_noanswer := con.QueryContext(ctx, sql_select_noanswer, username)
	helpers.ErrorCheck(err_noanswer)
	for row_noanswer.Next() {
		var (
			phone_db, nama_db, source_db, statuscrmsales_dua_db, notecrmsales_db string
			updatecrmsales_db, updatedatecrmsales_db                             string
		)

		err_noanswer = row_noanswer.Scan(&phone_db, &nama_db, &source_db, &statuscrmsales_dua_db, &notecrmsales_db,
			&updatecrmsales_db, &updatedatecrmsales_db)

		helpers.ErrorCheck(err_noanswer)

		obj_listnoanswer.Crmnoanswer_phone = phone_db
		obj_listnoanswer.Crmnoanswer_nama = nama_db
		obj_listnoanswer.Crmnoanswer_source = source_db
		obj_listnoanswer.Crmnoanswer_tipe = statuscrmsales_dua_db
		obj_listnoanswer.Crmnoanswer_note = notecrmsales_db
		obj_listnoanswer.Crmnoanswer_update = updatecrmsales_db + ", " + updatedatecrmsales_db
		arraobj_listnoanswer = append(arraobj_listnoanswer, obj_listnoanswer)
		msg = "Success"
	}
	defer row_noanswer.Close()

	//LIST INVALID
	var obj_listinvalid entities.Model_crmmemberlistinvalid
	var arraobj_listinvalid []entities.Model_crmmemberlistinvalid
	sql_select_invalid := ""
	sql_select_invalid += "SELECT "
	sql_select_invalid += "B.phone , B.nama, B.source,  "
	sql_select_invalid += "A.updatecrmsales, to_char(COALESCE(A.updatedatecrmsales,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select_invalid += "FROM " + configs.DB_tbl_trx_crmsales + " as A "
	sql_select_invalid += "JOIN " + configs.DB_tbl_trx_usersales + " as B ON B.phone = A.phone "
	sql_select_invalid += "WHERE A.username=$1  "
	sql_select_invalid += "AND A.statuscrmsales_satu='INVALID'   "
	if startdate != "" {
		sql_select_invalid += "AND A.updatedatecrmsales >='" + tglnow_start.Format("YYYY-MM-DD") + " 00:00:00' "
		sql_select_invalid += "AND A.updatedatecrmsales <='" + tglnow_end.Format("YYYY-MM-DD") + " 23:59:59' "
		sql_select_invalid += "ORDER BY A.updatedatecrmsales DESC   "
	} else {
		sql_select_invalid += "ORDER BY A.updatedatecrmsales DESC LIMIT 70  "
	}

	row_invalid, err_invalid := con.QueryContext(ctx, sql_select_invalid, username)
	helpers.ErrorCheck(err_invalid)
	for row_invalid.Next() {
		var (
			phone_db, nama_db, source_db             string
			updatecrmsales_db, updatedatecrmsales_db string
		)

		err_invalid = row_invalid.Scan(&phone_db, &nama_db, &source_db,
			&updatecrmsales_db, &updatedatecrmsales_db)

		helpers.ErrorCheck(err_invalid)

		obj_listinvalid.Crminvalid_phone = phone_db
		obj_listinvalid.Crminvalid_nama = nama_db
		obj_listinvalid.Crminvalid_source = source_db
		obj_listinvalid.Crminvalid_update = updatecrmsales_db + ", " + updatedatecrmsales_db
		arraobj_listinvalid = append(arraobj_listinvalid, obj_listinvalid)
		msg = "Success"
	}
	defer row_invalid.Close()

	obj.Sales_deposit = _GetSalesStatus(username, "DEPOSIT", "")
	obj.Sales_depositsum = float32(_GetSalesStatus(username, "DEPOSIT", "SUM"))
	obj.Sales_reject = _GetSalesStatus(username, "REJECT", "")
	obj.Sales_noanswer = _GetSalesStatus(username, "NOANSWER", "")
	obj.Sales_invalid = _GetSalesStatus(username, "INVALID", "")
	obj.Sales_listdeposit = arraobj_listdeposit
	obj.Sales_listnoanswer = arraobj_listnoanswer
	obj.Sales_listinvalid = arraobj_listinvalid
	arraobj = append(arraobj, obj)

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_employee(admin, username, password, iddepart, name, phone, status, sData, idrecord string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_employee, "username", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_employee + ` (
					username , password, iddepartement, 
					nmemployee , phoneemployee, statusemployee, 
					createemployee, createdateemployee
				) values (
					$1, $2, $3,
					$4, $5, $6,
					$7, $8
				)
			`
			hashpass := helpers.HashPasswordMD5(password)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_employee, "INSERT",
				username, hashpass, iddepart, name, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
				log.Println(msg_insert)
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		if password != "" {
			sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_employee + `  
				SET password=$1, 
				iddepartement=$2, nmemployee=$3,  phoneemployee=$4, statusemployee=$5, 
				updateemployee=$6, updatedateemployee=$7  
				WHERE username=$8 
			`
			hashpass := helpers.HashPasswordMD5(password)
			flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_employee, "UPDATE",
				hashpass, iddepart, name, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

			if flag_update {
				msg = "Succes"
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		} else {
			sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_employee + `  
				SET iddepartement=$1, nmemployee=$2, phoneemployee=$3, statusemployee=$4, 
				updateemployee=$5, updatedateemployee=$6  
				WHERE username=$7  
			`
			flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_employee, "UPDATE",
				iddepart, name, phone, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

			if flag_update {
				msg = "Succes"
				log.Println(msg_update)
			} else {
				log.Println(msg_update)
			}
		}

	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func _GetSalesStatus(username string, status, tipe string) int {
	con := db.CreateCon()
	ctx := context.Background()

	var (
		total       int
		total_float float64
	)

	sql_select := ""
	sql_select += "SELECT "
	if tipe == "SUM" {
		sql_select += "sum(deposit) as total "
	} else {
		sql_select += "count(idcrmsales) as total "
	}

	sql_select += "FROM " + configs.DB_tbl_trx_crmsales + " "
	sql_select += "WHERE username = $1 "
	if status == "INVALID" {
		sql_select += "AND statuscrmsales_satu = '" + status + "' "
	} else {
		sql_select += "AND statuscrmsales_dua = '" + status + "' "
	}
	log.Println(sql_select)
	row := con.QueryRowContext(ctx, sql_select, username)
	if tipe == "SUM" {
		switch e := row.Scan(&total_float); e {
		case sql.ErrNoRows:
		case nil:
		default:
			helpers.ErrorCheck(e)
		}
		total = int(total_float)
	} else {
		switch e := row.Scan(&total); e {
		case sql.ErrNoRows:
		case nil:
		default:
			helpers.ErrorCheck(e)
		}
	}
	return total
}
