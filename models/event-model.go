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
		idevent , idwebagen, nmevent,  
		startevent , endevent, 
		createevent, to_char(COALESCE(createdateevent,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updateevent, to_char(COALESCE(updatedateevent,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_trx_event + `  
		ORDER BY createdateevent DESC     
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idevent_db, idwebagen_db                                               int
			nmevent_db, startevent_db, endevent_db                                 string
			createevent_db, createdateevent_db, updateevent_db, updatedateevent_db string
		)

		err = row.Scan(&idevent_db, &idwebagen_db,
			&nmevent_db, &startevent_db, &endevent_db,
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
		obj.Event_name = nmevent_db
		obj.Event_startevent = startevent_db
		obj.Event_endevent = endevent_db
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
	admin, nmproviderslot, image, slug, title, descp, status, sData string,
	display, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_mst_providerslot + ` (
					idproviderslot , nmproviderslot, providerslot_display,  
					providerslot_counter , providerslot_status, providerslot_image,  
					providerslot_slug , providerslot_title, providerslot_descp,   
					createproviderslot, createdateproviderslot
				) values (
					$1, $2, $3, 
					$4, $5, $6, 
					$7, $8, $9, 
					$10, $11 
				)
			`
		field_column := configs.DB_tbl_mst_providerslot + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_providerslot, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), nmproviderslot,
			display, 0, status, image, slug, title, descp,
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
				` + configs.DB_tbl_mst_providerslot + `  
				SET nmproviderslot =$1, providerslot_display=$2, 
				providerslot_status=$3, providerslot_image=$4,
				providerslot_slug=$5, providerslot_title=$6,
				providerslot_descp=$7, 
				updateproviderslot=$8, updatedateproviderslot=$9  
				WHERE idproviderslot=$10 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_domain, "UPDATE",
			nmproviderslot, display, status, image, slug, title, descp,
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