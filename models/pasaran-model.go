package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_pasaranHome() (helpers.Response, error) {
	var obj entities.Model_pasaran
	var arraobj []entities.Model_pasaran
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idpasarantogel , nmpasarantogel, 
			urlpasaran , pasarandiundi, jamjadwal::text, displaypasaran, statuspasaran, 
			createpasarantogel, COALESCE(createdatepasarantogel,now()), updatepasarantogel, COALESCE(updatedatepasarantogel,now())  
			FROM ` + configs.DB_tbl_mst_pasaran + ` 
			ORDER BY displaypasaran ASC  
		`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			displaypasaran_db                                                                                     int
			idpasarantogel_db, nmpasarantogel_db, urlpasaran_db, pasarandiundi_db, jamjadwal_db, statuspasaran_db string
			createpasarantogel_db, createdatepasarantogel_db, updatepasarantogel_db, updatedatepasarantogel_db    string
		)

		err = row.Scan(
			&idpasarantogel_db, &nmpasarantogel_db, &urlpasaran_db, &pasarandiundi_db, &jamjadwal_db, &displaypasaran_db, &statuspasaran_db,
			&createpasarantogel_db, &createdatepasarantogel_db, &updatepasarantogel_db, &updatedatepasarantogel_db)

		helpers.ErrorCheck(err)
		statuspasaran := "HIDE"
		statuspasarancss := configs.STATUS_CANCEL
		create := ""
		update := ""
		if createpasarantogel_db != "" {
			create = createpasarantogel_db + ", " + createdatepasarantogel_db
		}
		if updatepasarantogel_db != "" {
			update = updatepasarantogel_db + ", " + updatedatepasarantogel_db
		}
		if statuspasaran_db == "Y" {
			statuspasaran = "SHOW"
			statuspasarancss = configs.STATUS_RUNNING
		}

		var (
			hasil_db string
		)
		sql_selectpasaran := `SELECT 
			Concat(datekeluaran , ' - ', nomorkeluaran) as hasil
			FROM ` + configs.DB_tbl_trx_keluaran + ` 
			WHERE idpasarantogel = $1 
			ORDER BY datekeluaran DESC LIMIT 1
		`
		row_keluaran := con.QueryRowContext(ctx, sql_selectpasaran, idpasarantogel_db)
		switch e_keluaran := row_keluaran.Scan(&hasil_db); e_keluaran {
		case sql.ErrNoRows:
		case nil:
		default:
			helpers.ErrorCheck(e_keluaran)
		}

		var (
			hasil_prediksi_db string
		)
		sql_selectprediksi := `SELECT concat( 
			dateprediksi , ' - ',  bbfsprediksi, ' - ', nomorprediksi ) as hasil_prediksi
			FROM ` + configs.DB_tbl_trx_prediksi + ` 
			WHERE idpasarantogel = $1 
			ORDER BY dateprediksi DESC LIMIT 1
		`
		row_prediksi := con.QueryRowContext(ctx, sql_selectprediksi, idpasarantogel_db)
		switch e_prediksi := row_prediksi.Scan(&hasil_prediksi_db); e_prediksi {
		case sql.ErrNoRows:
		case nil:
		default:
			helpers.ErrorCheck(e_prediksi)
		}

		obj.Pasaran_id = idpasarantogel_db
		obj.Pasaran_name = nmpasarantogel_db
		obj.Pasaran_url = urlpasaran_db
		obj.Pasaran_diundi = pasarandiundi_db
		obj.Pasaran_jamjadwal = jamjadwal_db
		obj.Pasaran_display = displaypasaran_db
		obj.Pasaran_status = statuspasaran
		obj.Pasaran_statuscss = statuspasarancss
		obj.Pasaran_keluaran = hasil_db
		obj.Pasaran_prediksi = hasil_prediksi_db
		obj.Pasaran_create = create
		obj.Pasaran_update = update
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
func Save_pasaran(admin, idrecord, nmpasarantogel, urlpasaran, pasarandiundi, jamjadwal, status, sData string, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_pasaran, "idpasarantogel", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_pasaran + ` (
					idpasarantogel , nmpasarantogel, urlpasaran, pasarandiundi, jamjadwal, displaypasaran, statuspasaran, 
					createpasarantogel, createdatepasarantogel
				) values (
					$1 ,$2, $3, $4, $5, $6, $7,
					$8, $9
				)
			`
			stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
			helpers.ErrorCheck(e_insert)
			defer stmt_insert.Close()
			res_newrecord, e_newrecord := stmt_insert.ExecContext(
				ctx,
				strings.ToUpper(idrecord), nmpasarantogel, urlpasaran, pasarandiundi, jamjadwal, display, status,
				admin,
				tglnow.Format("YYYY-MM-DD HH:mm:ss"))
			helpers.ErrorCheck(e_newrecord)
			insert, e := res_newrecord.RowsAffected()
			helpers.ErrorCheck(e)
			if insert > 0 {
				flag = true
				msg = "Succes"
				log.Println("Data Berhasil di save")
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_pasaran + `  
				SET nmpasarantogel=$1,urlpasaran=$2,pasarandiundi=$3, jamjadwal=$4, displaypasaran=$5, statuspasaran=$6,
				updatepasarantogel=$7, updatedatepasarantogel=$8 
				WHERE idpasarantogel =$9 
			`
		stmt_record, e := con.PrepareContext(ctx, sql_update)
		helpers.ErrorCheck(e)
		rec_record, e_record := stmt_record.ExecContext(
			ctx,
			nmpasarantogel, urlpasaran, pasarandiundi, jamjadwal, display, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			idrecord)
		helpers.ErrorCheck(e_record)
		update_record, e_record := rec_record.RowsAffected()
		helpers.ErrorCheck(e_record)

		defer stmt_record.Close()
		if update_record > 0 {
			flag = true
			msg = "Succes"
			log.Printf("Update PASARAN Success : %s\n", idrecord)
		} else {
			log.Println("Update PASARAN failed")
		}
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Fetch_keluaran(idpasaran string) (helpers.Response, error) {
	var obj entities.Model_keluaran
	var arraobj []entities.Model_keluaran
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idtrxkeluaran , datekeluaran, 
			periodekeluaran , nomorkeluaran
			FROM ` + configs.DB_tbl_trx_keluaran + ` 
			WHERE idpasarantogel=$1 
			ORDER BY datekeluaran DESC LIMIT 365 
		`

	row, err := con.QueryContext(ctx, sql_select, idpasaran)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idtrxkeluaran_db, periodekeluaran_db int
			datekeluaran_db, nomorkeluaran_db    string
		)

		err = row.Scan(
			&idtrxkeluaran_db, &datekeluaran_db, &periodekeluaran_db, &nomorkeluaran_db)

		helpers.ErrorCheck(err)

		obj.Keluaran_id = idtrxkeluaran_db
		obj.Keluaran_tanggal = datekeluaran_db
		obj.Keluaran_periode = idpasaran + "-" + strconv.Itoa(periodekeluaran_db)
		obj.Keluaran_nomor = nomorkeluaran_db
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
func Save_keluaran(admin, idpasaran, tanggal, nomor string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_trx_keluaran, "datekeluaran", tanggal, "idpasarantogel", idpasaran)
	if !flag {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_keluaran + ` (
				idtrxkeluaran , idpasarantogel, datekeluaran, periodekeluaran, nomorkeluaran, 
				createkeluaran, createdatekeluaran
			) values (
				$1 ,$2, $3, $4, $5,
				$6, $7
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_trx_keluaran + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		field_column_periode := strings.TrimSpace(idpasaran) + "-" + tglnow.Format("YYYY")
		idperiode_counter := Get_counter(field_column_periode)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			idpasaran, tanggal, idperiode_counter, nomor,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))
		helpers.ErrorCheck(e_newrecord)
		insert, e := res_newrecord.RowsAffected()
		helpers.ErrorCheck(e)
		if insert > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di save")
		}
	} else {
		msg = "Duplicate Entry"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Delete_keluaran(admin, idpasaran string, idtrxkeluaran int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_keluaran, "idtrxkeluaran", strconv.Itoa(idtrxkeluaran))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_keluaran + ` 
			WHERE idtrxkeluaran=$1 AND idpasarantogel=$2 
		`
		stmt_delete, e_delete := con.PrepareContext(ctx, sql_delete)
		helpers.ErrorCheck(e_delete)
		defer stmt_delete.Close()
		rec_delete, e_delete := stmt_delete.ExecContext(ctx, idtrxkeluaran, idpasaran)

		helpers.ErrorCheck(e_delete)
		delete, e := rec_delete.RowsAffected()
		helpers.ErrorCheck(e)
		if delete > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di delete")
		}
	} else {
		msg = "Data Not Found"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Fetch_prediksi(idpasaran string) (helpers.Response, error) {
	var obj entities.Model_prediksi
	var arraobj []entities.Model_prediksi
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idprediksi, dateprediksi, 
			bbfsprediksi , nomorprediksi
			FROM ` + configs.DB_tbl_trx_prediksi + ` 
			WHERE idpasarantogel=$1 
			ORDER BY dateprediksi DESC LIMIT 181  
		`

	row, err := con.QueryContext(ctx, sql_select, idpasaran)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idprediksi_db                                      int
			dateprediksi_db, bbfsprediksi_db, nomorprediksi_db string
		)

		err = row.Scan(&idprediksi_db, &dateprediksi_db, &bbfsprediksi_db, &nomorprediksi_db)

		helpers.ErrorCheck(err)

		obj.Prediksi_id = idprediksi_db
		obj.Prediksi_tanggal = dateprediksi_db
		obj.Prediksi_bbfs = bbfsprediksi_db
		obj.Prediksi_nomor = nomorprediksi_db
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
func Save_prediksi(admin, idpasaran, tanggal, bbfs, nomor string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_trx_prediksi, "dateprediksi", tanggal, "idpasarantogel", idpasaran)
	if !flag {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_prediksi + ` (
				idprediksi  , idpasarantogel, dateprediksi, bbfsprediksi, nomorprediksi, 
				createprediksi, createdateprediksi
			) values (
				$1 ,$2, $3, $4, $5,
				$6, $7
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_trx_prediksi + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			idpasaran, tanggal, bbfs, nomor,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))
		helpers.ErrorCheck(e_newrecord)
		insert, e := res_newrecord.RowsAffected()
		helpers.ErrorCheck(e)
		if insert > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di save")
		}
	} else {
		msg = "Duplicate Entry"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
func Delete_prediksi(admin, idpasaran string, idprediksi int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_prediksi, "idprediksi ", strconv.Itoa(idprediksi))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_prediksi + ` 
			WHERE idprediksi =$1 AND idpasarantogel=$2 
		`
		stmt_delete, e_delete := con.PrepareContext(ctx, sql_delete)
		helpers.ErrorCheck(e_delete)
		defer stmt_delete.Close()
		rec_delete, e_delete := stmt_delete.ExecContext(ctx, idprediksi, idpasaran)

		helpers.ErrorCheck(e_delete)
		delete, e := rec_delete.RowsAffected()
		helpers.ErrorCheck(e)
		if delete > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di delete")
		}
	} else {
		msg = "Data Not Found"
	}

	if flag {
		res.Status = fiber.StatusOK
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	} else {
		res.Status = fiber.StatusBadRequest
		res.Message = msg
		res.Record = nil
		res.Time = time.Since(render_page).String()
	}

	return res, nil
}
