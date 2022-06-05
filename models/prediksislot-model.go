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

func Fetch_prediksislotHome(idproviderslot int) (helpers.Response, error) {
	var obj entities.Model_prediksislot
	var arraobj []entities.Model_prediksislot
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "A.idgameslot , B.nmproviderslot, A.nmgameslot, A.gameslot_prediksi, "
	sql_select += "A.gameslot_image , A.gameslot_status,   "
	sql_select += "A.creategameslot, to_char(COALESCE(A.createdategameslot,now()), 'YYYY-MM-DD HH24:MI:SS'), "
	sql_select += "A.updategameslot, to_char(COALESCE(A.updatedategameslot,now()), 'YYYY-MM-DD HH24:MI:SS') "
	sql_select += "FROM " + configs.DB_tbl_trx_gameslot + " as A  "
	sql_select += "JOIN " + configs.DB_tbl_mst_providerslot + " as B ON B.idproviderslot = A.idproviderslot  "
	if idproviderslot > 0 {
		sql_select += "WHERE A.idproviderslot = '" + strconv.Itoa(idproviderslot) + "' "
		sql_select += "ORDER BY A.gameslot_prediksi DESC "
	} else {
		sql_select += "ORDER BY A.gameslot_prediksi DESC "
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idgameslot_db, gameslot_prediksi_db                                                int
			nmproviderslot_db, nmgameslot_db, gameslot_image_db, gameslot_status_db            string
			creategameslot_db, createdategameslot_db, updategameslot_db, updatedategameslot_db string
		)

		err = row.Scan(&idgameslot_db, &nmproviderslot_db, &nmgameslot_db,
			&gameslot_prediksi_db, &gameslot_image_db, &gameslot_status_db,
			&creategameslot_db, &createdategameslot_db, &updategameslot_db, &updatedategameslot_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if creategameslot_db != "" {
			create = creategameslot_db + ", " + createdategameslot_db
		}
		if updategameslot_db != "" {
			update = updategameslot_db + ", " + updatedategameslot_db
		}

		obj.Prediksislot_id = idgameslot_db
		obj.Prediksislot_nmprovider = nmproviderslot_db
		obj.Prediksislot_name = nmgameslot_db
		obj.Prediksislot_prediksi = gameslot_prediksi_db
		obj.Prediksislot_image = gameslot_image_db
		obj.Prediksislot_status = gameslot_status_db
		obj.Prediksislot_create = create
		obj.Prediksislot_update = update
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
func Save_prediksislot(
	admin, nmgameslot, image, status, sData string,
	idproviderslot, prediksi, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_trx_gameslot + ` (
					idgameslot , idproviderslot, nmgameslot,  
					gameslot_prediksi , gameslot_image, gameslot_status,  
					creategameslot, createdategameslot
				) values (
					$1, $2, $3, 
					$4, $5, $6, 
					$7, $8 
				)
			`
		field_column := configs.DB_tbl_trx_gameslot + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_gameslot, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), idproviderslot, nmgameslot,
			prediksi, image, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_trx_gameslot + `  
				SET idproviderslot=$1, nmgameslot=$2, 
				gameslot_prediksi=$3, gameslot_image=$4, gameslot_status=$5, 
				updategameslot=$6, updatedategameslot=$7  
				WHERE idgameslot=$8 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_domain, "UPDATE",
			idproviderslot, nmgameslot, prediksi, image, status,
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
func Delete_prediksislot(admin string, idrecord, idproviderslot int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()

	sql_delete := `
		DELETE FROM
		` + configs.DB_tbl_trx_gameslot + ` 
		WHERE idgameslot=$1 AND idproviderslot=$2 
	`

	flag_episode, msg_episode := Exec_SQL(sql_delete, configs.DB_tbl_trx_gameslot, "DELETE", idrecord, idproviderslot)

	if flag_episode {
		msg = "Succes"
		log.Println(msg_episode)
	} else {
		log.Println(msg_episode)
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()
	return res, nil
}
