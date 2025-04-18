package login

import (
	"net/http"

	"imageresizerservice/app/ctx/appCtx"
	"imageresizerservice/app/users/login/sendLink"
	"imageresizerservice/app/users/login/useLink"
)

func Router(mux *http.ServeMux, ac *appCtx.AppCtx) {
	useLink.Router(mux, ac)
}

func RouterLoggedOut(mux *http.ServeMux, ac *appCtx.AppCtx) {
	sendLink.Router(mux, ac)
	useLink.Router(mux, ac)
}
