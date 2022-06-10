package models

import (
	"context"
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
