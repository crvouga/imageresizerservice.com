package projects

import (
	"imageresizerservice/app/ctx/appCtx"
	"imageresizerservice/app/projects/createProject"
	"imageresizerservice/app/projects/deleteProject"
	"imageresizerservice/app/projects/editProject"
	"imageresizerservice/app/projects/listProjects"
	"imageresizerservice/app/projects/projectPage"

	"net/http"
)

func Router(mux *http.ServeMux, ac *appCtx.AppCtx) {
	createProject.Router(mux, ac)
	editProject.Router(mux, ac)
	deleteProject.Router(mux, ac)
	listProjects.Router(mux, ac)
	projectPage.Router(mux, ac)
}
