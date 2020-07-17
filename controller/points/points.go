package points

import (
	"bytes"
	"fmt"

	"github.com/wmentor/gencset/db"
	"github.com/wmentor/serv"
)

func init() {
	serv.Register("GET", "/points", handler)
}

type Loc struct {
	Id        int64
	Code      string
	Latitude  float64
	Longitude float64
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

	sth, err := dbh.Prepare("SELECT id, code, latitude, longitude FROM locs WHERE id > $1 AND NOT skip ORDER BY id LIMIT 100")
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
			if err := rows.Scan(&l.Id, &l.Code, &l.Latitude, &l.Longitude); err != nil {
				panic(err)
			}

			id = l.Id

			maker.WriteString(l.Code)
			maker.WriteRune(';')
			maker.WriteString(fmt.Sprintf("%.9f;%.9f\n", l.Latitude, l.Longitude))

			cnt++
		}

		c.Write(maker.Bytes())

		if cnt < 100 {
			break
		}
	}
}
