package main

import (
	"image-resizer-service/deps"
	"image-resizer-service/email/send_email"
	"image-resizer-service/login"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User-agent: *\nAllow: /"))
	})

	d := deps.Deps{
		SendEmail: &send_email.FakeSendEmail{},
	}

	mux.HandleFunc("/login", login.RespondLoginPage)
	mux.HandleFunc("/login/send-link", login.RespondSendLink(&d))
	mux.HandleFunc("/login/sent-link", login.RespondSentLinkPage)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	port := "8080"

	log.Printf("Server live here http://localhost:%s/ \n", port)

	http.ListenAndServe(":8080", mux)

}
