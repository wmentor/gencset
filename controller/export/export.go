package export

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/wmentor/gencset/db"
	"github.com/wmentor/serv"
)

func init() {
	serv.Register("GET", "/export", handler)
}

type Loc struct {
	Id        int64
	Name      string
	Code      string
	Latitude  float64
	Longitude float64
	Forms     string
}

func handler(c *serv.Context) {

	c.SetContentType("text/plain; charset=utf-8")
	c.WriteHeader(200)

	id := int64(0)

	maker := bytes.NewBuffer(nil)

	dbh, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer dbh.Close()

	sth, err := dbh.Prepare("SELECT id, name, code, latitude, longitude, forms FROM locs WHERE id > $1 AND NOT skip ORDER BY id LIMIT 100")
	if err != nil {
		panic(err)
	}
	defer sth.Close()

	var l Loc

	for {

		maker.Reset()

		rows, err := sth.Query(id)
		if err != nil {
			panic(err)
		}

		cnt := 0

		for rows.Next() {
			if err := rows.Scan(&l.Id, &l.Name, &l.Code, &l.Latitude, &l.Longitude, &l.Forms); err != nil {
				panic(err)
			}

			id = l.Id

			for _, f := range strings.Split(l.Forms, "/") {
				str := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(f)), "ё", "е")

				maker.WriteString(str)
				maker.WriteRune(';')
				maker.WriteString(l.Code)
				maker.WriteRune(';')
				maker.WriteString(fmt.Sprintf("%.9f;%.9f\n", l.Latitude, l.Longitude))
			}

			cnt++
		}

		c.Write(maker.Bytes())

		if cnt < 100 {
			break
		}
	}
}
