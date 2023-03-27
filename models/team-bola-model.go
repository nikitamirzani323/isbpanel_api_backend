package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_teamHome() (helpers.Response, error) {
	var obj entities.Model_Team
	var arraobj []entities.Model_Team
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
		idteam , nmteam, imgteam,  statusteam,  
		createteam, to_char(COALESCE(createdateteam,now()), 'YYYY-MM-DD HH24:MI:SS'), 
		updateteam, to_char(COALESCE(updatedateteam,now()), 'YYYY-MM-DD HH24:MI:SS') 
		FROM ` + configs.DB_tbl_mst_bola_team + `  
		ORDER BY updatedateteam DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idteam_db                                                          int
			nmteam_db, imgteam_db, statusteam_db                               string
			createteam_db, createdateteam_db, updateteam_db, updatedateteam_db string
		)

		err = row.Scan(&idteam_db, &nmteam_db, &imgteam_db, &statusteam_db,
			&createteam_db, &createdateteam_db, &updateteam_db, &updatedateteam_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createteam_db != "" {
			create = createteam_db + ", " + createdateteam_db
		}
		if updateteam_db != "" {
			update = updateteam_db + ", " + updatedateteam_db
		}

		obj.Team_id = idteam_db
		obj.Team_name = nmteam_db
		obj.Team_image = imgteam_db
		obj.Team_status = statusteam_db
		obj.Team_create = create
		obj.Team_update = update
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
func Save_team(admin, name, image, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_bola_team, "nmteam", name)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_bola_team + ` (
					idleague , nmleague, imgleague, statusleague, 
					createleague, createdateleague
				) values (
					$1, $2, $3, $4, 
					$5, $6
				)
			`
			field_column := configs.DB_tbl_mst_bola_league + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_bola_league, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, image, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_bola_league + `  
				SET nmleague=$1, imgleague=$2, statusleague=$3, 
				updateleague=$4, updatedateleague=$5  
				WHERE idleague=$6 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_domain, "UPDATE",
			name, image, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
