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

func Fetch_departementHome() (helpers.Response, error) {
	var obj entities.Model_departement
	var arraobj []entities.Model_departement
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			iddepartement , nmdepartement,  
			createdepartement, to_char(COALESCE(createdatedepartement,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatedepartement, to_char(COALESCE(updatedatedepartement,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_departement + `  
			ORDER BY createdatedepartement DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddepartement_db, nmdepartement_db                                                             string
			createdepartement_db, createdatedepartement_db, updatedepartement_db, updatedatedepartement_db string
		)

		err = row.Scan(&iddepartement_db, &nmdepartement_db,
			&createdepartement_db, &createdatedepartement_db, &updatedepartement_db, &updatedatedepartement_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createdepartement_db != "" {
			create = createdepartement_db + ", " + createdatedepartement_db
		}
		if updatedepartement_db != "" {
			update = updatedepartement_db + ", " + updatedatedepartement_db
		}

		obj.Departement_id = iddepartement_db
		obj.Departement_name = nmdepartement_db
		obj.Departement_create = create
		obj.Departement_update = update
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
func Save_departement(admin, nmdepartement, sData, idrecord string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_departement, "iddepartement", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_departement + ` (
					iddepartement , nmdepartement, 
					createdepartement, createdatedepartement
				) values (
					$1, $2, 
					$3, $4
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_departement, "INSERT",
				idrecord, nmdepartement,
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
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_departement + `  
				SET nmdepartement =$1, 
				updatedepartement=$2, updatedatedepartement=$3 
				WHERE iddepartement=$4 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_departement, "UPDATE",
			nmdepartement,
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
