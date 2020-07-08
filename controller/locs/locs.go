package locs

import (
	"strconv"
	"strings"

	"github.com/wmentor/gencset/db"
	"github.com/wmentor/serv"
)

func init() {

	serv.Register("GET", "/", handlePage)
	serv.Register("GET", "/locs", handlePage)
	serv.Register("GET", "/locs/:id", handlePage)
	serv.Register("GET", "/locs/:id/add", handleEditPage)
	serv.Register("POST", "/locs/:id/add", handleSavePage)
	serv.Register("GET", "/locs/:id/edit/:loc_id", handleEditPage)
	serv.Register("POST", "/locs/:id/edit/:loc_id", handleSavePage)

}

type Loc struct {
	Id        int64
	Name      string
	Code      string
	Latitude  float64
	Longitude float64
	Forms     string
	ParentId  int64
}

func handlePage(c *serv.Context) {

	c.SetContentType("text/html; charset=utf-8")

	id := c.ParamInt("id")

	vars := map[string]interface{}{}

	dbh, err := db.Get()
	if err != nil || dbh == nil {
		panic("db connect broken")
	}
	defer dbh.Close()

	if id == 0 {
		vars["title"] = "Geo locations"
		vars["id"] = int64(0)
	} else {
		vars["id"] = id

		row := dbh.QueryRow("SELECT id, name, code, latitude, longitude, forms, parent_id FROM locs WHERE id = $1", id)
		if row == nil {
			c.WriteRedirect("/locs")
			return
		}

		l := &Loc{}

		if err := row.Scan(&l.Id, &l.Name, &l.Code, &l.Latitude, &l.Longitude, &l.Forms, &l.ParentId); err != nil {
			panic(err)
		}

		vars["current"] = l
	}

	rows, err := dbh.Query("SELECT id, name, code, latitude, longitude, forms, parent_id "+
		"FROM locs WHERE parent_id=$1 ORDER BY LOWER(name)", id)
	if err != nil {
		panic(err)
	}

	var list []*Loc

	for rows.Next() {
		l := &Loc{}

		if err := rows.Scan(&l.Id, &l.Name, &l.Code, &l.Latitude, &l.Longitude, &l.Forms, &l.ParentId); err != nil {
			panic(err)
		}

		list = append(list, l)
	}

	vars["locs"] = list

	c.WriteHeader(200)
	c.Render("locs.tt", vars)
}

func handleEditPage(c *serv.Context) {

	id := c.ParamInt64("loc_id")
	parent_id := c.ParamInt64("id")

	var l Loc

	if id == 0 {
		l.ParentId = parent_id
	} else {

		dbh, err := db.Get()
		if err != nil {
			panic(err)
		}
		defer dbh.Close()

		row := dbh.QueryRow("SELECT id, name, code, forms, latitude, longitude, parent_id FROM locs WHERE id = $1", id)
		if row == nil {
			panic("no row")
		}

		err = row.Scan(&l.Id, &l.Name, &l.Code, &l.Forms, &l.Latitude, &l.Longitude, &l.ParentId)
		if err != nil {
			panic(err)
		}
	}

	vars := map[string]interface{}{
		"loc": l,
	}

	c.SetContentType("text/html; charset=utf-8")
	c.WriteHeader(200)
	c.Render("locs/edit.tt", vars)
}

func handleSavePage(c *serv.Context) {

	var l Loc

	id := c.ParamInt64("loc_id")
	parent_id := c.ParamInt64("id")

	l.Name = c.FormValue("name")

	coords := c.FormValue("coords")
	cv := strings.Split(coords, ",")

	if len(cv) != 2 {
		c.WriteRedirect("/locs/" + strconv.FormatInt(parent_id, 10))
		return
	}

	l.Latitude, _ = strconv.ParseFloat(strings.TrimSpace(cv[0]), 64)
	l.Longitude, _ = strconv.ParseFloat(strings.TrimSpace(cv[1]), 64)

	l.Forms = c.FormValue("forms")

	dbh, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer dbh.Close()

	tx, err := dbh.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	if id == 0 {
		if _, err := tx.Exec("INSERT INTO locs(name,latitude,longitude,forms,parent_id) VALUES ($1,$2,$3,$4,$5)",
			l.Name, l.Latitude, l.Longitude, l.Forms, parent_id); err == nil {
			tx.Commit()
		}
	} else {
		if _, err := tx.Exec("UPDATE locs SET name=$1, latitude=$2, longitude = $3, forms=$4 WHERE id=$5",
			l.Name, l.Latitude, l.Longitude, l.Forms, id); err == nil {
			tx.Commit()
		} else {
			panic(err)
		}
	}

	c.WriteRedirect("/locs/" + strconv.FormatInt(parent_id, 10))
}
