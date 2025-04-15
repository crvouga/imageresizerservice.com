package sendLinkPage

import (
	"net/http"
	"net/url"

	"imageresizerservice/app/ui/page"
	"imageresizerservice/app/users/loginWithEmailLink/routes"
	"imageresizerservice/library/static"
)

type Data struct {
	Action     string
	EmailError string
	Email      string
	JsPath     string
}

func Router(mux *http.ServeMux) {
	mux.HandleFunc(routes.SendLinkPage, Respond())
}

func Respond() http.HandlerFunc {
	htmlPath := static.GetSiblingPath("sendLinkPage.html")

	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{
			Action:     routes.SendLinkAction,
			Email:      r.URL.Query().Get("Email"),
			EmailError: r.URL.Query().Get("ErrorEmail"),
		}

		page.Respond(htmlPath, data)(w, r)
	}
}

type RedirectErrorArgs struct {
	Email      string
	EmailError string
}

func RedirectError(w http.ResponseWriter, r *http.Request, args RedirectErrorArgs) {
	u, _ := url.Parse(routes.SendLinkPage)
	q := u.Query()
	q.Set("Email", args.Email)
	q.Set("ErrorEmail", args.EmailError)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.SendLinkPage, http.StatusSeeOther)
}
