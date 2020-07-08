package stat

import (
	"fmt"

	"github.com/wmentor/gencset/db"
	"github.com/wmentor/serv"
)

func init() {
	serv.Register("GET", "/stat", handler)
}

func handler(c *serv.Context) {

	dbh, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer dbh.Close()

	if row := dbh.QueryRow("SELECT COUNT(id) FROM locs"); err == nil {
		var cnt int64

		if err := row.Scan(&cnt); err == nil {
			c.SetContentType("text/plain; charset=utf-8")
			c.WriteHeader(200)
			c.WriteString(fmt.Sprintf("Locs=%d\n", cnt))
			return
		}
	} else {
		panic(err)
	}

}
