package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"imageresizerservice/deps"
	"imageresizerservice/email/emailOutbox"
	"imageresizerservice/email/sendEmail"
	"imageresizerservice/keyValueDb"
	"imageresizerservice/static"
	"imageresizerservice/uow"
	"imageresizerservice/users"
	"imageresizerservice/users/loginWithEmailLink/link/linkDb"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatalf("Failed to open in-memory SQLite database: %v", err)
	}

	defer db.Close()

	keyValueDbHashMap := keyValueDb.ImplHashMap{}

	d := deps.Deps{
		SendEmail:   &sendEmail.ImplFake{},
		LinkDb:      &linkDb.ImplKeyValueDb{Db: &keyValueDbHashMap},
		UowFactory:  uow.UowFactory{Db: db},
		KeyValueDb:  &keyValueDbHashMap,
		EmailOutbox: &emailOutbox.ImplKeyValueDb{Db: &keyValueDbHashMap},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	users.Router(mux, &d)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := static.ServeStaticAssets(w, r)
		if err != nil {
			http.Redirect(w, r, "/login-with-email-link/login-page", http.StatusSeeOther)
		}
	})

	port := "8080"

	log.Printf("Server live here http://localhost:%s/ \n", port)

	http.ListenAndServe(":8080", mux)

	// emailOutboxWorker.Start(&d, 10*time.Second)
}
