package index

import (
	"github.com/wmentor/serv"
)

func init() {

	serv.Register("GET", "/", handler)
	serv.Register("GET", "index.html", handler)

}

func handler(c *serv.Context) {

	c.SetContentType("text/html; charset=utf-8")
	c.WriteHeader(200)
	c.Render("index.tt", nil)

}
