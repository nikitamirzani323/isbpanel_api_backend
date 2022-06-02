package models

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_gameHome(search string, page int) (helpers.Responsemovie, error) {
	var obj entities.Model_game
	var arraobj []entities.Model_game
	var res helpers.Responsemovie
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 50
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idgame) as totalgame  "
	sql_selectcount += "FROM " + configs.DB_tbl_mst_game + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(nmgame) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nmgame) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idgame ,nmgame, "
	sql_select += "creategame, to_char(COALESCE(createdategame,now()), 'YYYY-MM-DD HH24:ii:ss') as createdategame, updategame, to_char(COALESCE(updatedategame,now()), 'YYYY-MM-DD HH24:ii:ss') as updatedategame "
	sql_select += "FROM " + configs.DB_tbl_mst_game + " "
	if search == "" {
		sql_select += "ORDER BY createdategame DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(nmgame) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(nmgame) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdategame DESC  LIMIT " + strconv.Itoa(perpage)
	}
	log.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idgame_db                                                          int
			nmgame_db                                                          string
			creategame_db, createdategame_db, updategame_db, updatedategame_db string
		)

		err = row.Scan(
			&idgame_db, &nmgame_db,
			&creategame_db, &createdategame_db, &updategame_db, &updatedategame_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if creategame_db != "" {
			create = creategame_db + ", " + createdategame_db
		}
		if updategame_db != "" {
			update = updategame_db + ", " + updatedategame_db
		}

		obj.Game_id = idgame_db
		obj.Game_name = nmgame_db
		obj.Game_create = create
		obj.Game_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_game(admin, sdata, name string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_game + ` (
				idgame , nmgame,  
				creategame, createdategame  
			) values (
				$1 ,$2, 
				$3, $4  
			)
		`
		field_column := configs.DB_tbl_mst_game + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_game, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name,
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
			` + configs.DB_tbl_mst_game + ` 
			SET nmgame=$1,   
			updategame=$2, updatedategame=$3 
			WHERE idgame=$4  
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_game, "UPDATE",
			name, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

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
