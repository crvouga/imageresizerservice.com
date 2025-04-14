package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"imageresizerservice/deps"
	"imageresizerservice/static"
	"imageresizerservice/users"
	"imageresizerservice/users/loginWithEmailLink/routes"
)

func main() {
	db := newDb()

	defer db.Close()

	d := deps.New(db)

	mux := http.NewServeMux()

	Router(mux, &d)

	addr := ":8080"

	log.Printf("Server live here http://localhost%s/ \n", addr)

	http.ListenAndServe(addr, mux)
}

func newDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	return db
}

func Router(mux *http.ServeMux, d *deps.Deps) {

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	users.Router(mux, d)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := static.ServeStaticAssets(w, r)
		if err != nil {
			http.Redirect(w, r, routes.SendLinkPage, http.StatusSeeOther)
		}
	})
}
