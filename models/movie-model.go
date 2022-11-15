package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/isbtotogroup/isbpanel_api_backend/configs"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/db"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/entities"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/helpers"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/nleeper/goment"
)

func Fetch_movieHome(search string, page, enable int) (helpers.Responsemovie, error) {
	var obj entities.Model_movie
	var arraobj []entities.Model_movie
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
	sql_selectcount += "COUNT(movieid) as totalmovie  "
	sql_selectcount += "FROM " + configs.DB_VIEW_MOVIE + "  "
	sql_selectcount += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
	if search != "" {
		sql_selectcount += "AND LOWER(movietitle) LIKE '%" + strings.ToLower(search) + "%' "
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
	sql_select += "movieid , movietitle, description, movietype, "
	sql_select += "rating , imdb, year, views, enabled, COALESCE(posted_id,0),label, urlthumbnail, slug,   "
	sql_select += "createmovie, to_char(COALESCE(createdatemovie,now()), 'YYYY-MM-DD HH24:MI:SS') , updatemovie, to_char(COALESCE(updatedatemovie,now()), 'YYYY-MM-DD HH24:MI:SS')  "
	sql_select += "FROM " + configs.DB_VIEW_MOVIE + "  "
	if search == "" {
		sql_select += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
		sql_select += "ORDER BY createdatemovie DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
		sql_select += "AND LOWER(movietitle) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdatemovie DESC LIMIT " + strconv.Itoa(perpage)
	}
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			movieid_db, year_db, views_db, enabled_db, posted_id_db                         int
			rating_db, imdb_db                                                              float32
			movietitle_db, movietype_db, description_db, label_db, urlthumbnail_db, slug_db string
			createmovie_db, createdatemovie_db, updatemovie_db, updatedatemovie_db          string
		)

		err = row.Scan(
			&movieid_db, &movietitle_db, &description_db, &movietype_db, &rating_db, &imdb_db, &year_db, &views_db,
			&enabled_db, &posted_id_db, &label_db, &urlthumbnail_db, &slug_db,
			&createmovie_db, &createdatemovie_db, &updatemovie_db, &updatedatemovie_db)

		helpers.ErrorCheck(err)
		status := "HIDE"
		statuscss := configs.STATUS_CANCEL
		create := ""
		update := ""
		createdt, _ := goment.New(createdatemovie_db)
		updatedt, _ := goment.New(updatedatemovie_db)

		if createmovie_db != "" {
			create = fmt.Sprintf("%s, %s", createmovie_db, createdt.Format("YYYY-MM-DD HH:mm:ss"))
		}
		if updatemovie_db != "" {
			update = fmt.Sprintf("%s, %s", updatemovie_db, updatedt.Format("YYYY-MM-DD HH:mm:ss"))
		}
		if enabled_db > 0 {
			status = "SHOW"
			statuscss = configs.STATUS_RUNNING
		}
		//GENRE
		var objmoviegenre entities.Model_moviegenre
		var arraobjmoviegenre []entities.Model_moviegenre
		sql_selectmoviegenre := `SELECT 
			A.idgenre, B.nmgenre 
			FROM ` + configs.DB_tbl_trx_moviegenre + ` as A 
			JOIN ` + configs.DB_tbl_mst_moviegenre + ` as B ON B.idgenre = A.idgenre 
			WHERE A.movieid = $1   
		`
		row_moviegenre, err := con.QueryContext(ctx, sql_selectmoviegenre, movieid_db)
		helpers.ErrorCheck(err)
		for row_moviegenre.Next() {
			var (
				idgenre_db int
				nmgenre_db string
			)
			err = row_moviegenre.Scan(&idgenre_db, &nmgenre_db)
			objmoviegenre.Moviegenre_id = idgenre_db
			objmoviegenre.Moviegenre_name = nmgenre_db
			arraobjmoviegenre = append(arraobjmoviegenre, objmoviegenre)
		}
		//SOURCE
		var objmoviesource entities.Model_moviesource
		var arraobjmoviesource []entities.Model_moviesource
		if movietype_db == "movie" {
			sql_selectmoviesource := `SELECT 
				id, url 
				FROM ` + configs.DB_tbl_mst_moviesource + ` 
				WHERE poster_id = $1  
			`
			row_moviesource, err := con.QueryContext(ctx, sql_selectmoviesource, movieid_db)
			helpers.ErrorCheck(err)
			nosource := 0
			for row_moviesource.Next() {
				nosource = nosource + 1
				var (
					id_db  int
					url_db string
				)
				err = row_moviesource.Scan(&id_db, &url_db)
				objmoviesource.Moviesource_id = id_db
				objmoviesource.Moviesource_stream = "STREAM-" + strconv.Itoa(nosource)
				objmoviesource.Moviesource_url = url_db
				arraobjmoviesource = append(arraobjmoviesource, objmoviesource)
			}
		}

		path_image := ""
		if urlthumbnail_db == "" {
			poster_image, poster_extension := _GetMedia(posted_id_db)
			path_image = "https://duniafilm.b-cdn.net/uploads/cache/poster_thumb/uploads/" + poster_extension + "/" + poster_image
		} else {
			path_image = urlthumbnail_db
		}
		imgcdn := "N"
		if urlthumbnail_db != "" {
			imgcdn = "Y"
		}

		obj.Movie_id = movieid_db
		obj.Movie_date = createdt.Format("YYYY-MM-DD HH:mm:ss")
		obj.Movie_type = movietype_db
		obj.Movie_title = movietitle_db
		obj.Movie_label = label_db
		obj.Movie_slug = slug_db
		obj.Movie_descp = description_db
		obj.Movie_imgcdn = imgcdn
		obj.Movie_thumbnail = path_image
		obj.Movie_year = year_db
		obj.Movie_rating = rating_db
		obj.Movie_imdb = imdb_db
		obj.Movie_view = views_db
		obj.Movie_comment = _GetTotalComment(movieid_db)
		obj.Movie_genre = arraobjmoviegenre
		obj.Movie_source = arraobjmoviesource
		obj.Movie_status = status
		obj.Movie_statuscss = statuscss
		obj.Movie_create = create
		obj.Movie_update = update
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
func Fetch_movieHomeNotCDN() (helpers.Response, error) {
	var obj entities.Model_movienotcdn
	var arraobj []entities.Model_movienotcdn
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "movieid , movietitle "
	sql_select += "FROM " + configs.DB_VIEW_MOVIE + "  "
	sql_select += "WHERE urlthumbnail = '' "
	sql_select += "ORDER BY createdatemovie DESC  "
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			movieid_db    int
			movietitle_db string
		)

		err = row.Scan(&movieid_db, &movietitle_db)

		helpers.ErrorCheck(err)

		obj.Movie_id = movieid_db
		obj.Movie_title = movietitle_db

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
func Fetch_movieHomeBanner() (helpers.Response, error) {
	var obj entities.Model_moviebanner
	var arraobj []entities.Model_moviebanner
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "moviebannerid, nmmoviebanner , urlimgmoviebanner, urldestinationmoviebanner, devicemoviebanner, "
	sql_select += "displaymoviebanner, statusmoviebanner "
	sql_select += "FROM " + configs.DB_tbl_trx_movie_banner + "  "
	sql_select += "ORDER BY displaymoviebanner ASC  "
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			moviebannerid_db, displaymoviebanner_db                                                                          int
			nmmoviebanner_db, urlimgmoviebanner_db, urldestinationmoviebanner_db, devicemoviebanner_db, statusmoviebanner_db string
		)

		err = row.Scan(&moviebannerid_db, &nmmoviebanner_db, &urlimgmoviebanner_db, &urldestinationmoviebanner_db,
			&devicemoviebanner_db, &displaymoviebanner_db, &statusmoviebanner_db)

		helpers.ErrorCheck(err)

		statusbanner := "HIDE"
		statusbannercss := configs.STATUS_CANCEL
		if statusmoviebanner_db == "Y" {
			statusbanner = "SHOW"
			statusbannercss = configs.STATUS_RUNNING
		}
		obj.Moviebanner_id = moviebannerid_db
		obj.Moviebanner_title = nmmoviebanner_db
		obj.Moviebanner_urlimage = urlimgmoviebanner_db
		obj.Moviebanner_urldestination = urldestinationmoviebanner_db
		obj.Moviebanner_device = devicemoviebanner_db
		obj.Moviebanner_display = displaymoviebanner_db
		obj.Moviebanner_status = statusbanner
		obj.Moviebanner_statuscss = statusbannercss

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
func Save_movie(admin, name, label, slug, tipemovie, descp, urlthum, listgenre, listsource, sdata string, idrecord, year, status int, imdb float32) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	temp_idrecord := ""
	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_movie + ` (
				movieid , movietitle, label, movietype, description, imdb, year, slug, enabled, urlthumbnail,    
				createmovie, createdatemovie
			) values (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10, 
				$11, $12
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_trx_movie + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		temp_idrecord = tglnow.Format("YY") + tglnow.Format("MM") + strconv.Itoa(idrecord_counter)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			temp_idrecord,
			name, label, tipemovie, descp, imdb, year, slug, status, urlthum,
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
		if flag {
			//GENRE
			json := []byte(listgenre)
			jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_genre_id, _ := jsonparser.GetInt(value, "movie_genre_id")
				movie_genre_name, _ := jsonparser.GetString(value, "movie_genre_name")

				log.Printf("%d - %s", int(movie_genre_id), movie_genre_name)
				sql_insertgenre := `
				insert into
					` + configs.DB_tbl_trx_moviegenre + ` (
						idmoviegenre , movieid, idgenre
					) values (
						$1,$2,$3
					)
				`
				stmt_insertgenre, e_insertgenre := con.PrepareContext(ctx, sql_insertgenre)
				helpers.ErrorCheck(e_insertgenre)
				defer stmt_insertgenre.Close()
				field_column := configs.DB_tbl_trx_moviegenre + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordgenre, e_newrecordgenre := stmt_insertgenre.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					temp_idrecord, movie_genre_id)
				helpers.ErrorCheck(e_newrecordgenre)
				insertgenre, e := res_newrecordgenre.RowsAffected()
				helpers.ErrorCheck(e)
				if insertgenre > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data GENERE Berhasil di save")
				}
			})

			//SOURCE
			jsonsource := []byte(listsource)
			jsonparser.ArrayEach(jsonsource, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_source_id, _ := jsonparser.GetInt(value, "movie_source_id")
				movie_source_name, _ := jsonparser.GetString(value, "movie_source_name")

				log.Printf("%d - %s", int(movie_source_id), movie_source_name)
				sql_insertsource := `
				insert into
					` + configs.DB_tbl_mst_moviesource + ` (
						id , poster_id, type, url
					) values (
						$1,$2,$3,$4
					)
				`
				stmt_insertsource, e_insertsource := con.PrepareContext(ctx, sql_insertsource)
				helpers.ErrorCheck(e_insertsource)
				defer stmt_insertsource.Close()
				field_column := configs.DB_tbl_mst_moviesource + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordsource, e_newrecordsource := stmt_insertsource.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					temp_idrecord, "embed", movie_source_name)
				helpers.ErrorCheck(e_newrecordsource)
				insertsource, e := res_newrecordsource.RowsAffected()
				helpers.ErrorCheck(e)
				if insertsource > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data SOURCE Berhasil di save")
				}
			})
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_trx_movie + ` 
			SET movietitle=$1, description=$2, slug=$3, urlthumbnail=$4, enabled=$5 ,
			updatemovie=$6, updatedatemovie=$7 
			WHERE movieid=$8 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_movie, "UPDATE",
			name, descp, slug, urlthum, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}

		if flag {
			//DELETE GENRE
			stmt_genre_delete, e_genre_delete := con.PrepareContext(ctx, `
					DELETE FROM  
					`+configs.DB_tbl_trx_moviegenre+`   
					WHERE movieid=$1  
			`)

			helpers.ErrorCheck(e_genre_delete)
			rec_genre_delete, e_rec_genre_delete := stmt_genre_delete.ExecContext(ctx, idrecord)
			helpers.ErrorCheck(e_rec_genre_delete)

			affect_genre_delete, err_affer_genre_delete := rec_genre_delete.RowsAffected()
			helpers.ErrorCheck(err_affer_genre_delete)

			defer stmt_genre_delete.Close()
			if affect_genre_delete > 0 {
				log.Printf("Delete genre : %d\n", idrecord)
			} else {
				log.Printf("Delete genre : %d Failed\n ", idrecord)
			}

			//DELETE SOURCE
			stmt_source_delete, e_source_delete := con.PrepareContext(ctx, `
					DELETE FROM  
					`+configs.DB_tbl_mst_moviesource+`   
					WHERE poster_id=$1  
			`)

			helpers.ErrorCheck(e_source_delete)
			rec_source_delete, e_rec_source_delete := stmt_source_delete.ExecContext(ctx, idrecord)
			helpers.ErrorCheck(e_rec_source_delete)

			affect_source_delete, err_affer_source_delete := rec_source_delete.RowsAffected()
			helpers.ErrorCheck(err_affer_source_delete)

			defer stmt_source_delete.Close()
			if affect_source_delete > 0 {
				log.Printf("Delete source : %d\n", idrecord)
			} else {
				log.Printf("Delete source : %d Failed\n ", idrecord)
			}

			//GENRE
			json := []byte(listgenre)
			jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_genre_id, _ := jsonparser.GetInt(value, "movie_genre_id")
				movie_genre_name, _ := jsonparser.GetString(value, "movie_genre_name")

				log.Printf("%d - %s", int(movie_genre_id), movie_genre_name)
				sql_insertgenre := `
					insert into
						` + configs.DB_tbl_trx_moviegenre + ` (
							idmoviegenre , movieid, idgenre
						) values (
							$1,$2,$3
						)
					`
				stmt_insertgenre, e_insertgenre := con.PrepareContext(ctx, sql_insertgenre)
				helpers.ErrorCheck(e_insertgenre)
				defer stmt_insertgenre.Close()
				field_column := configs.DB_tbl_trx_moviegenre + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordgenre, e_newrecordgenre := stmt_insertgenre.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					idrecord, movie_genre_id)
				helpers.ErrorCheck(e_newrecordgenre)
				insertgenre, e := res_newrecordgenre.RowsAffected()
				helpers.ErrorCheck(e)
				if insertgenre > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data GENERE Berhasil di save")
				}
			})

			//SOURCE
			jsonsource := []byte(listsource)
			jsonparser.ArrayEach(jsonsource, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_source_id, _ := jsonparser.GetInt(value, "movie_source_id")
				movie_source_name, _ := jsonparser.GetString(value, "movie_source_name")

				log.Printf("%d - %s", int(movie_source_id), movie_source_name)
				sql_insertsource := `
				insert into
					` + configs.DB_tbl_mst_moviesource + ` (
						id , poster_id, type, url
					) values (
						$1,$2,$3,$4
					)
				`
				stmt_insertsource, e_insertsource := con.PrepareContext(ctx, sql_insertsource)
				helpers.ErrorCheck(e_insertsource)
				defer stmt_insertsource.Close()
				field_column := configs.DB_tbl_mst_moviesource + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordsource, e_newrecordsource := stmt_insertsource.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					idrecord, "embed", movie_source_name)
				helpers.ErrorCheck(e_newrecordsource)
				insertsource, e := res_newrecordsource.RowsAffected()
				helpers.ErrorCheck(e)
				if insertsource > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data SOURCE Berhasil di save")
				}
			})
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
func Delete_movie(admin string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_movie, "movieid", strconv.Itoa(idrecord))

	if flag {
		//DELETE GENRE
		sql_deletegenre := `
			DELETE FROM
			` + configs.DB_tbl_trx_moviegenre + ` 
			WHERE movieid=$1 
		`
		stmt_deletegenre, e_deletegenre := con.PrepareContext(ctx, sql_deletegenre)
		helpers.ErrorCheck(e_deletegenre)
		defer stmt_deletegenre.Close()
		rec_deletegenre, e_deletegenre := stmt_deletegenre.ExecContext(ctx, idrecord)

		helpers.ErrorCheck(e_deletegenre)
		deletegenre, egenre := rec_deletegenre.RowsAffected()
		helpers.ErrorCheck(egenre)
		if deletegenre > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Genre Berhasil di delete")
		}

		//DELETE SOURCE
		sql_deletesource := `
			DELETE FROM
			` + configs.DB_tbl_mst_moviesource + ` 
			WHERE poster_id=$1 
		`
		stmt_deletesource, e_deletesource := con.PrepareContext(ctx, sql_deletesource)
		helpers.ErrorCheck(e_deletesource)
		defer stmt_deletesource.Close()
		rec_deletesource, e_deletesource := stmt_deletesource.ExecContext(ctx, idrecord)

		helpers.ErrorCheck(e_deletesource)
		deletesource, esource := rec_deletesource.RowsAffected()
		helpers.ErrorCheck(esource)
		if deletesource > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Source Berhasil di delete")
		}

		//DELETE MOVIE
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_movie + ` 
			WHERE movieid=$1 
		`
		stmt_delete, e_delete := con.PrepareContext(ctx, sql_delete)
		helpers.ErrorCheck(e_delete)
		defer stmt_delete.Close()
		rec_delete, e_delete := stmt_delete.ExecContext(ctx, idrecord)

		helpers.ErrorCheck(e_delete)
		delete, e := rec_delete.RowsAffected()
		helpers.ErrorCheck(e)
		if delete > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Movie Berhasil di delete")
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
func Fetch_movieseriesHome(search string, page, enable int) (helpers.Responsemovie, error) {
	var obj entities.Model_movieseries
	var arraobj []entities.Model_movieseries
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
	sql_selectcount += "COUNT(movieid) as totalmovie  "
	sql_selectcount += "FROM " + configs.DB_VIEW_MOVIESERIES + "  "
	sql_selectcount += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
	if search != "" {
		sql_selectcount += "AND LOWER(movietitle) LIKE '%" + strings.ToLower(search) + "%' "
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
	sql_select += "movieid , movietitle, description, movietype, "
	sql_select += "rating , imdb, year, views, enabled, COALESCE(posted_id,0),label, urlthumbnail, slug,   "
	sql_select += "createmovie, to_char(COALESCE(createdatemovie,now()), 'YYYY-MM-DD HH24:MI:SS'), updatemovie, to_char(COALESCE(updatedatemovie,now()), 'YYYY-MM-DD HH24:MI:SS') "
	sql_select += "FROM " + configs.DB_VIEW_MOVIESERIES + "  "
	if search == "" {
		sql_select += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
		sql_select += "ORDER BY createdatemovie DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE enabled = '" + strconv.Itoa(enable) + "' "
		sql_select += "AND LOWER(movietitle) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdatemovie DESC LIMIT " + strconv.Itoa(perpage)
	}
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			movieid_db, year_db, views_db, enabled_db, posted_id_db                         int
			rating_db, imdb_db                                                              float32
			movietitle_db, movietype_db, description_db, label_db, urlthumbnail_db, slug_db string
			createmovie_db, createdatemovie_db, updatemovie_db, updatedatemovie_db          string
		)

		err = row.Scan(
			&movieid_db, &movietitle_db, &description_db, &movietype_db, &rating_db, &imdb_db, &year_db, &views_db,
			&enabled_db, &posted_id_db, &label_db, &urlthumbnail_db, &slug_db,
			&createmovie_db, &createdatemovie_db, &updatemovie_db, &updatedatemovie_db)

		helpers.ErrorCheck(err)
		status := "HIDE"
		statuscss := configs.STATUS_CANCEL
		create := ""
		update := ""
		createdt, _ := goment.New(createdatemovie_db)
		updatedt, _ := goment.New(updatedatemovie_db)

		if createmovie_db != "" {
			create = fmt.Sprintf("%s, %s", createmovie_db, createdt.Format("YYYY-MM-DD HH:mm:ss"))
		}
		if updatemovie_db != "" {
			update = fmt.Sprintf("%s, %s", updatemovie_db, updatedt.Format("YYYY-MM-DD HH:mm:ss"))
		}
		if enabled_db > 0 {
			status = "SHOW"
			statuscss = configs.STATUS_RUNNING
		}
		//GENRE
		var objmoviegenre entities.Model_moviegenre
		var arraobjmoviegenre []entities.Model_moviegenre
		sql_selectmoviegenre := `SELECT 
			A.idgenre, B.nmgenre 
			FROM ` + configs.DB_tbl_trx_moviegenre + ` as A 
			JOIN ` + configs.DB_tbl_mst_moviegenre + ` as B ON B.idgenre = A.idgenre 
			WHERE A.movieid = $1   
		`
		row_moviegenre, err := con.QueryContext(ctx, sql_selectmoviegenre, movieid_db)
		helpers.ErrorCheck(err)
		for row_moviegenre.Next() {
			var (
				idgenre_db int
				nmgenre_db string
			)
			err = row_moviegenre.Scan(&idgenre_db, &nmgenre_db)
			objmoviegenre.Moviegenre_id = idgenre_db
			objmoviegenre.Moviegenre_name = nmgenre_db
			arraobjmoviegenre = append(arraobjmoviegenre, objmoviegenre)
		}

		//SEASON
		var objmovieseason entities.Model_movieseason
		var arraobjmovieseason []entities.Model_movieseason
		sql_selectmovieseason := `SELECT 
			id, title, position 
			FROM ` + configs.DB_tbl_mst_season + ` 
			WHERE poster_id = $1   
			ORDER BY position ASC 
		`
		row_movieseason, err := con.QueryContext(ctx, sql_selectmovieseason, movieid_db)
		helpers.ErrorCheck(err)
		for row_movieseason.Next() {
			var (
				id_db, position_db int
				title_db           string
			)
			err = row_movieseason.Scan(&id_db, &title_db, &position_db)
			totalepisode := _GetTotalEpisode(id_db)
			objmovieseason.Movieseason_id = id_db
			objmovieseason.Movieseason_name = title_db
			objmovieseason.Movieseason_episodetotal = totalepisode
			objmovieseason.Movieseason_display = position_db
			arraobjmovieseason = append(arraobjmovieseason, objmovieseason)
		}

		path_image := ""
		if urlthumbnail_db == "" {
			poster_image, poster_extension := _GetMedia(posted_id_db)
			path_image = "https://duniafilm.b-cdn.net/uploads/cache/poster_thumb/uploads/" + poster_extension + "/" + poster_image
		} else {
			path_image = urlthumbnail_db
		}

		imgcdn := "N"
		if urlthumbnail_db != "" {
			imgcdn = "Y"
		}
		obj.Movie_id = movieid_db
		obj.Movie_date = createdt.Format("YYYY-MM-DD HH:mm:ss")
		obj.Movie_type = movietype_db
		obj.Movie_title = movietitle_db
		obj.Movie_label = label_db
		obj.Movie_slug = slug_db
		obj.Movie_descp = description_db
		obj.Movie_imgcdn = imgcdn
		obj.Movie_thumbnail = path_image
		obj.Movie_year = year_db
		obj.Movie_rating = rating_db
		obj.Movie_imdb = imdb_db
		obj.Movie_view = views_db
		obj.Movie_genre = arraobjmoviegenre
		obj.Movie_season = arraobjmovieseason
		obj.Movie_status = status
		obj.Movie_statuscss = statuscss
		obj.Movie_create = create
		obj.Movie_update = update
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
func Save_movieseries(admin, name, label, slug, tipemovie, descp, urlthum, listgenre, sdata string, idrecord, year, status int, imdb float32) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	con := db.CreateCon()
	ctx := context.Background()
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	temp_idrecord := ""
	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_movie + ` (
				movieid , movietitle, label, movietype, description, imdb, year, slug, enabled, urlthumbnail,    
				createmovie, createdatemovie
			) values (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10, 
				$11, $12
			)
		`
		stmt_insert, e_insert := con.PrepareContext(ctx, sql_insert)
		helpers.ErrorCheck(e_insert)
		defer stmt_insert.Close()
		field_column := configs.DB_tbl_trx_movie + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		temp_idrecord = tglnow.Format("YY") + tglnow.Format("MM") + strconv.Itoa(idrecord_counter)
		res_newrecord, e_newrecord := stmt_insert.ExecContext(
			ctx,
			temp_idrecord,
			name, label, tipemovie, descp, imdb, year, slug, status, urlthum,
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
		if flag {
			//GENRE
			json := []byte(listgenre)
			jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_genre_id, _ := jsonparser.GetInt(value, "movie_genre_id")
				movie_genre_name, _ := jsonparser.GetString(value, "movie_genre_name")

				log.Printf("%d - %s", int(movie_genre_id), movie_genre_name)
				sql_insertgenre := `
				insert into
					` + configs.DB_tbl_trx_moviegenre + ` (
						idmoviegenre , movieid, idgenre
					) values (
						$1,$2,$3
					)
				`
				stmt_insertgenre, e_insertgenre := con.PrepareContext(ctx, sql_insertgenre)
				helpers.ErrorCheck(e_insertgenre)
				defer stmt_insertgenre.Close()
				field_column := configs.DB_tbl_trx_moviegenre + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordgenre, e_newrecordgenre := stmt_insertgenre.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					temp_idrecord, movie_genre_id)
				helpers.ErrorCheck(e_newrecordgenre)
				insertgenre, e := res_newrecordgenre.RowsAffected()
				helpers.ErrorCheck(e)
				if insertgenre > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data GENERE Berhasil di save")
				}
			})
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_trx_movie + ` 
			SET movietitle=$1, description=$2, slug=$3, urlthumbnail=$4, 
			updatemovie=$5, updatedatemovie=$6 
			WHERE movieid=$7 
		`
		stmt_update, e_update := con.PrepareContext(ctx, sql_update)
		helpers.ErrorCheck(e_update)
		defer stmt_update.Close()
		res_newrecord, e_newrecord := stmt_update.ExecContext(
			ctx,
			name, descp, slug, urlthum,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)
		helpers.ErrorCheck(e_newrecord)
		update, e := res_newrecord.RowsAffected()
		helpers.ErrorCheck(e)
		if update > 0 {
			flag = true
			msg = "Succes"
			log.Println("Data Berhasil di update")
		}
		if flag {
			//DELETE GENRE
			stmt_genre_delete, e_genre_delete := con.PrepareContext(ctx, `
					DELETE FROM  
					`+configs.DB_tbl_trx_moviegenre+`   
					WHERE movieid=$1  
			`)

			helpers.ErrorCheck(e_genre_delete)
			rec_genre_delete, e_rec_genre_delete := stmt_genre_delete.ExecContext(ctx, idrecord)
			helpers.ErrorCheck(e_rec_genre_delete)

			affect_genre_delete, err_affer_genre_delete := rec_genre_delete.RowsAffected()
			helpers.ErrorCheck(err_affer_genre_delete)

			defer stmt_genre_delete.Close()
			if affect_genre_delete > 0 {
				log.Printf("Delete genre : %d\n", idrecord)
			} else {
				log.Printf("Delete genre : %d Failed\n ", idrecord)
			}

			//GENRE
			json := []byte(listgenre)
			jsonparser.ArrayEach(json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				movie_genre_id, _ := jsonparser.GetInt(value, "movie_genre_id")
				movie_genre_name, _ := jsonparser.GetString(value, "movie_genre_name")

				log.Printf("%d - %s", int(movie_genre_id), movie_genre_name)
				sql_insertgenre := `
					insert into
						` + configs.DB_tbl_trx_moviegenre + ` (
							idmoviegenre , movieid, idgenre
						) values (
							$1,$2,$3
						)
					`
				stmt_insertgenre, e_insertgenre := con.PrepareContext(ctx, sql_insertgenre)
				helpers.ErrorCheck(e_insertgenre)
				defer stmt_insertgenre.Close()
				field_column := configs.DB_tbl_trx_moviegenre + tglnow.Format("YYYY")
				idrecord_counter := Get_counter(field_column)
				res_newrecordgenre, e_newrecordgenre := stmt_insertgenre.ExecContext(
					ctx,
					tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecord_counter),
					idrecord, movie_genre_id)
				helpers.ErrorCheck(e_newrecordgenre)
				insertgenre, e := res_newrecordgenre.RowsAffected()
				helpers.ErrorCheck(e)
				if insertgenre > 0 {
					flag = true
					msg = "Succes"
					log.Println("Data GENERE Berhasil di save")
				}
			})
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
func Fetch_season(idmovie int) (helpers.Response, error) {
	var obj entities.Model_movieseason
	var arraobj []entities.Model_movieseason
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			id , title, position
			FROM ` + configs.DB_tbl_mst_season + ` 
			WHERE poster_id=$1 
			ORDER BY position ASC   
		`

	row, err := con.QueryContext(ctx, sql_select, idmovie)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			id_db, position_db int
			title_db           string
		)

		err = row.Scan(&id_db, &title_db, &position_db)
		helpers.ErrorCheck(err)
		totalepisode := _GetTotalEpisode(id_db)

		obj.Movieseason_id = id_db
		obj.Movieseason_name = title_db
		obj.Movieseason_episodetotal = totalepisode
		obj.Movieseason_display = position_db
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
func Save_season(admin, name, sdata string, idrecord, idmovie, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_season + ` (
				id ,poster_id, title, position
			) values (
				$1 ,$2, $3, $4
			)
		`
		field_column := configs.DB_tbl_mst_season + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_season, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			idmovie, name, display)

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_mst_season + ` 
			SET title=$1, position=$2 
			WHERE id=$3 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_season, "UPDATE", name, display, idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
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
func Delete_season(admin string, idrecord, idmovie int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_mst_season, "id", strconv.Itoa(idrecord), "poster_id", strconv.Itoa(idmovie))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_mst_season + ` 
			WHERE id=$1 AND poster_id=$2 
		`

		flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_mst_season, "DELETE", idrecord, idmovie)

		if flag_delete {
			flag = true
			msg = "Succes"
			log.Println(msg_delete)
		} else {
			log.Println(msg_delete)
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
func Fetch_episode(idseason int) (helpers.Response, error) {
	var obj entities.Model_movieepisode
	var arraobj []entities.Model_movieepisode
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			id , title, position
			FROM ` + configs.DB_tbl_mst_episode + ` 
			WHERE season_id=$1 
			ORDER BY position ASC   
		`

	row, err := con.QueryContext(ctx, sql_select, idseason)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			id_db, position_db int
			title_db           string
		)

		err = row.Scan(&id_db, &title_db, &position_db)
		helpers.ErrorCheck(err)

		//SOURCE
		var objmoviesource entities.Model_moviesource
		var arraobjmoviesource []entities.Model_moviesource
		sql_selectmoviesource := `SELECT 
			id, url 
			FROM ` + configs.DB_tbl_mst_moviesource + ` 
			WHERE episode_id = $1   
		`
		row_moviesource, err := con.QueryContext(ctx, sql_selectmoviesource, id_db)
		helpers.ErrorCheck(err)
		nosource := 0
		for row_moviesource.Next() {
			nosource = nosource + 1
			var (
				id_db  int
				url_db string
			)
			err = row_moviesource.Scan(&id_db, &url_db)
			objmoviesource.Moviesource_id = id_db
			objmoviesource.Moviesource_stream = "STREAM-" + strconv.Itoa(nosource)
			objmoviesource.Moviesource_url = url_db
			arraobjmoviesource = append(arraobjmoviesource, objmoviesource)
		}

		obj.Movieepisode_id = id_db
		obj.Movieepisode_seasonid = idseason
		obj.Movieepisode_name = title_db
		obj.Movieepisode_display = position_db
		obj.Movieepisode_source = arraobjmoviesource
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
func Save_episode(admin, name, urlsource, sdata string, idrecord, idseason, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_episode + ` (
				id ,season_id, title, position, enabled
			) values (
				$1 ,$2, $3, $4, $5 
			)
		`
		field_column := configs.DB_tbl_mst_episode + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		temp_idrecord := tglnow.Format("YY") + tglnow.Format("MM") + strconv.Itoa(idrecord_counter)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_episode, "INSERT",
			temp_idrecord,
			idseason, name, display, "1")

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}

		sql_insertsource := `
			insert into
				` + configs.DB_tbl_mst_moviesource + ` (
					id , episode_id, type, url
				) values (
					$1,$2,$3,$4
				)
		`
		field_columnsource := configs.DB_tbl_mst_moviesource + tglnow.Format("YYYY")
		idrecordsource_counter := Get_counter(field_columnsource)
		flag_insertsource, msg_insertsource := Exec_SQL(sql_insertsource, configs.DB_tbl_mst_moviesource, "INSERT",
			tglnow.Format("YY")+tglnow.Format("MM")+strconv.Itoa(idrecordsource_counter),
			temp_idrecord, "embed", urlsource)

		if flag_insertsource {
			flag = true
			msg = "Succes"
			log.Println(msg_insertsource)
		} else {
			log.Println(msg_insertsource)
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
func Delete_episode(admin string, idrecord, idseason int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false

	flag = CheckDBTwoField(configs.DB_tbl_mst_episode, "id", strconv.Itoa(idrecord), "season_id", strconv.Itoa(idseason))
	if flag {
		sql_deletesource := `
			DELETE FROM
			` + configs.DB_tbl_mst_moviesource + ` 
			WHERE episode_id=$1 
		`

		flag_deletesource, msg_deletesource := Exec_SQL(sql_deletesource, configs.DB_tbl_mst_moviesource, "DELETE", idrecord)

		if flag_deletesource {
			flag = true
			msg = "Succes"
			log.Println(msg_deletesource)

			sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_mst_episode + ` 
				WHERE id=$1 AND season_id=$2 
			`

			flag_episode, msg_episode := Exec_SQL(sql_delete, configs.DB_tbl_mst_episode, "DELETE", idrecord, idseason)

			if flag_episode {
				flag = true
				msg = "Succes"
				log.Println(msg_episode)
			} else {
				log.Println(msg_episode)
			}
		} else {
			log.Println(msg_deletesource)
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
func Save_moviebanner(admin, name, urlimg, urldestination, device, status, sdata string, idrecord, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_trx_movie_banner + ` (
				moviebannerid ,nmmoviebanner, 
				urlimgmoviebanner, urldestinationmoviebanner, devicemoviebanner, displaymoviebanner ,statusmoviebanner, 
				createmoviebanner, createdatemoviebanner
			) values (
				$1 ,$2, 
				$3, $4, $5, $6, $7, 
				$8, $9
			)
		`
		field_column := configs.DB_tbl_trx_movie_banner + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		temp_idrecord := tglnow.Format("YY") + tglnow.Format("MM") + strconv.Itoa(idrecord_counter)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_trx_movie_banner, "INSERT",
			temp_idrecord,
			name, urlimg, urldestination, device, display, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_trx_movie_banner + ` 
			SET nmmoviebanner=$1, urlimgmoviebanner=$2, urldestinationmoviebanner=$3, devicemoviebanner=$4, displaymoviebanner=$5 ,statusmoviebanner=$6 , 
			updatemoviebanner=$7, updatedatemoviebanner=$8 
			WHERE moviebannerid=$9 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_trx_movie_banner, "UPDATE",
			name, urlimg, urldestination, device, display, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
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
func Delete_moviebanner(admin string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false

	flag = CheckDB(configs.DB_tbl_trx_movie_banner, "moviebannerid", strconv.Itoa(idrecord))
	if flag {
		sql_delete := `
			DELETE FROM
			` + configs.DB_tbl_trx_movie_banner + ` 
			WHERE moviebannerid=$1 
		`

		flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_trx_movie_banner, "DELETE", idrecord)

		if flag_delete {
			flag = true
			msg = "Succes"
			log.Println(msg_delete)
		} else {
			log.Println(msg_delete)
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
func Fetch_movieminiHome(search string) (helpers.Response, error) {
	var obj entities.Model_minimovie
	var arraobj []entities.Model_minimovie
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := ""
	sql_select += "SELECT "
	sql_select += "movieid , movietitle, movietype "
	sql_select += "FROM " + configs.DB_VIEW_MOVIE + "  "
	if search == "" {
		sql_select += "ORDER BY createdatemovie DESC  LIMIT 100 "
	} else {
		sql_select += "WHERE LOWER(movietitle) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY createdatemovie DESC LIMIT 100 "
	}
	log.Println(sql_select)
	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			movieid_db                  int
			movietitle_db, movietype_db string
		)

		err = row.Scan(&movieid_db, &movietitle_db, &movietype_db)

		helpers.ErrorCheck(err)

		obj.Movie_id = movieid_db
		obj.Movie_title = movietitle_db
		obj.Movie_type = movietype_db
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

func Fetch_genre() (helpers.Response, error) {
	var obj entities.Model_genre
	var arraobj []entities.Model_genre
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idgenre , nmgenre, genredisplay, 
			creategenre, to_char(COALESCE(createdategenre,now()), 'YYYY-MM-DD HH24:MI:SS') , updategenre, to_char(COALESCE(updatedategenre,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_moviegenre + ` 
			ORDER BY genredisplay ASC   
		`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idgenre_db, genredisplay_db                                            int
			nmgenre_db                                                             string
			creategenre_db, createdategenre_db, updategenre_db, updatedategenre_db string
		)

		err = row.Scan(
			&idgenre_db, &nmgenre_db, &genredisplay_db,
			&creategenre_db, &createdategenre_db, &updategenre_db, &updatedategenre_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if creategenre_db != "" {
			create = creategenre_db + ", " + createdategenre_db
		}
		if updategenre_db != "" {
			update = updategenre_db + ", " + updatedategenre_db
		}

		obj.Genre_id = idgenre_db
		obj.Genre_name = nmgenre_db
		obj.Genre_display = genredisplay_db
		obj.Genre_create = create
		obj.Genre_update = update
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
func Save_genre(admin, name, sdata string, idrecord, display int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sdata == "New" {
		sql_insert := `
			insert into
			` + configs.DB_tbl_mst_moviegenre + ` (
				idgenre , nmgenre, genredisplay,  
				creategenre, createdategenre
			) values (
				$1 ,$2, $3, 
				$4, $5
			)
		`
		field_column := configs.DB_tbl_mst_moviegenre + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_moviegenre, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter),
			name, display,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			flag = true
			msg = "Succes"
			log.Println(msg_insert)
		} else {
			log.Println(msg_insert)
		}
	} else {
		sql_update := `
			UPDATE 
			` + configs.DB_tbl_mst_moviegenre + ` 
			SET nmgenre=$1, genredisplay=$2, 
			updategenre=$3, updatedategenre=$4 
			WHERE idgenre=$5 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_moviegenre, "UPDATE",
			name, display,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
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
func Delete_genre(admin string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	render_page := time.Now()
	flag := false
	flag_movie := false

	flag = CheckDB(configs.DB_tbl_mst_moviegenre, "idgenre", strconv.Itoa(idrecord))
	flag_movie = CheckDB(configs.DB_tbl_trx_moviegenre, "idgenre", strconv.Itoa(idrecord))
	if flag {
		if !flag_movie {
			sql_delete := `
				DELETE FROM
				` + configs.DB_tbl_mst_moviegenre + ` 
				WHERE idgenre=$1 
			`
			flag_delete, msg_delete := Exec_SQL(sql_delete, configs.DB_tbl_mst_moviegenre, "DELETE", idrecord)

			if flag_delete {
				flag = true
				msg = "Succes"
				log.Println(msg_delete)
			} else {
				log.Println(msg_delete)
			}
		} else {
			msg = "Cannot Delete"
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
func _GetMedia(idrecord int) (string, string) {
	con := db.CreateCon()
	ctx := context.Background()
	url := ""
	extension := ""

	sql_select := `SELECT
		url, extension   
		FROM ` + configs.DB_tbl_mst_mediatable + `  
		WHERE idmediatable = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&url, &extension); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return url, extension
}
func _GetTotalEpisode(idrecord int) int {
	con := db.CreateCon()
	ctx := context.Background()
	total := 0

	sql_select := `SELECT
		count(id) as total  
		FROM ` + configs.DB_tbl_mst_episode + `  
		WHERE season_id = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&total); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return total
}
func _GetTotalComment(idrecord int) int {
	con := db.CreateCon()
	ctx := context.Background()
	total := 0

	sql_select := `SELECT
		count(idcomment) as total  
		FROM ` + configs.DB_tbl_trx_comment + `  
		WHERE idposter = $1 
	`
	row := con.QueryRowContext(ctx, sql_select, idrecord)
	switch e := row.Scan(&total); e {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e)
	}
	return total
}
