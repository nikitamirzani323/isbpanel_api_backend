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

func Fetch_bannerHome() (helpers.Response, error) {
	var obj entities.Model_banner
	var arraobj []entities.Model_banner
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idbanner , nmbanner, urlbanner, urlwebsite, devicebanner, posisibanner, displaybanner, statusbanner, 
		createbanner, to_char(COALESCE(createdatebanner,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updatebanner, to_char(COALESCE(updatedatebanner,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_mst_banner + `  
		ORDER BY displaybanner ASC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idbanner_db, displaybanner_db                                                               int
			nmbanner_db, urlbanner_db, urlwebsite_db, devicebanner_db, posisibanner_db, statusbanner_db string
			createbanner_db, createdatebanner_db, updatebanner_db, updatedatebanner_db                  string
		)

		err = row.Scan(&idbanner_db, &nmbanner_db, &urlbanner_db, &urlwebsite_db, &devicebanner_db, &posisibanner_db, &displaybanner_db,
			&statusbanner_db, &createbanner_db, &createdatebanner_db, &updatebanner_db, &updatedatebanner_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createbanner_db != "" {
			create = createbanner_db + ", " + createdatebanner_db
		}
		if updatebanner_db != "" {
			update = updatebanner_db + ", " + updatedatebanner_db
		}

		obj.Banner_id = idbanner_db
		obj.Banner_name = nmbanner_db
		obj.Banner_url = urlbanner_db
		obj.Banner_urlwebsite = urlwebsite_db
		obj.Banner_device = devicebanner_db
		obj.Banner_posisi = posisibanner_db
		obj.Banner_display = displaybanner_db
		obj.Banner_status = statusbanner_db
		obj.Banner_create = create
		obj.Banner_update = update
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

func Duplicate_banner(admin string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()

	sql_select := `SELECT 
			nmbanner , urlbanner, urlwebsite, devicebanner, posisibanner, displaybanner, statusbanner  
			FROM ` + configs.DB_tbl_mst_banner + `  
			WHERE idbanner = $1 
	`

	row, err := con.QueryContext(ctx, sql_select, idrecord)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			displaybanner_db                                                                            int
			nmbanner_db, urlbanner_db, urlwebsite_db, devicebanner_db, posisibanner_db, statusbanner_db string
		)

		err = row.Scan(&nmbanner_db, &urlbanner_db, &urlwebsite_db,
			&devicebanner_db, &posisibanner_db, &displaybanner_db, &statusbanner_db)

		helpers.ErrorCheck(err)

		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_banner + ` (
				idbanner , nmbanner, urlbanner, urlwebsite, 
				devicebanner, posisibanner, displaybanner, statusbanner, 
				createbanner, createdatebanner
			) values (
				$1 ,$2, $3, $4,
				$5, $6, $7, $8,
				$9, $10  
			)
		`
		field_column := configs.DB_tbl_mst_banner + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_banner, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			nmbanner_db, urlbanner_db, urlwebsite_db, devicebanner_db, posisibanner_db, displaybanner_db, statusbanner_db,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
func Save_banner(admin, sdata, nmbanner, urlbanner, urlwebsite, devicebanner, posisibanner, status string, idrecord, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_banner + ` (
				idbanner , nmbanner, urlbanner, urlwebsite, 
				devicebanner, posisibanner, displaybanner, statusbanner, 
				createbanner, createdatebanner
			) values (
				$1 ,$2, $3, $4,
				$5, $6, $7, $8,
				$9, $10  
			)
		`
		field_column := configs.DB_tbl_mst_banner + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_banner, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			nmbanner, urlbanner, urlwebsite, devicebanner, posisibanner, display, status,
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
				` + configs.DB_tbl_mst_banner + `  
				SET nmbanner=$1, urlbanner=$2, urlwebsite=$3, devicebanner=$4,
				posisibanner=$5, displaybanner=$6, statusbanner=$7,
				updatebanner=$8, updatedatebanner=$9 
				WHERE idbanner=$10  
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_banner, "UPDATE",
			nmbanner, urlbanner, urlwebsite, devicebanner, posisibanner, display, status,
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
