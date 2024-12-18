package page

import (
	"html/template"
	"net/http"
)

func Respond(pageTemplatePath string, pageData any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./page/page.html", "./ui/icons.html", "./ui/header.html", pageTemplatePath)

		if err != nil {
			errStr := err.Error()
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, pageData); err != nil {
			errStr := err.Error()
			http.Error(w, errStr, http.StatusInternalServerError)
		}
	}
}
