package useLinkPage

import (
	"imageresizerservice/deps"
	"imageresizerservice/page"
	"imageresizerservice/static"
	"imageresizerservice/users/loginWithEmailLink/routes"
	"net/http"
	"net/url"
)

func Router(mux *http.ServeMux) {
	mux.HandleFunc(routes.UseLinkPage, Respond())
}

type Data struct {
	Action string
	LinkId string
}

func Respond() http.HandlerFunc {
	htmlPath := static.GetSiblingPath("useLinkPage.html")
	return func(w http.ResponseWriter, r *http.Request) {
		data := Data{
			Action: routes.UseLinkAction,
			LinkId: r.URL.Query().Get("linkId"),
		}

		page.Respond(htmlPath, data)(w, r)
	}
}

func ToUrl(d *deps.Deps, linkId string) string {
	path := ToPath(linkId)
	u, _ := url.Parse(d.BaseUrl + path)
	return u.String()
}

func ToPath(linkId string) string {
	u, _ := url.Parse(routes.UseLinkPage)
	q := u.Query()
	q.Set("linkId", linkId)
	u.RawQuery = q.Encode()
	return u.String()
}
