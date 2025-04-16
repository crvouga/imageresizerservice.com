package deleteProject

import (
	"imageresizerservice/app/ctx/appContext"
	"imageresizerservice/app/ctx/reqCtx"
	"imageresizerservice/app/projects/project"
	"imageresizerservice/app/projects/project/projectID"
	"imageresizerservice/app/projects/projectRoutes"
	"imageresizerservice/app/ui/page"
	"imageresizerservice/library/static"
	"net/http"
)

func Router(mux *http.ServeMux, appCtx *appContext.AppCtx) {
	mux.HandleFunc(projectRoutes.ProjectDelete, Respond(appCtx))
}

type Data struct {
	Project     *project.Project
	ProjectPage string
}

func Respond(appCtx *appContext.AppCtx) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			respondPost(appCtx, w, r)
		} else {
			respondGet(appCtx, w, r)
		}
	}
}

func respondGet(appCtx *appContext.AppCtx, w http.ResponseWriter, r *http.Request) {
	req := reqCtx.FromHttpRequest(appCtx, r)
	logger := req.Logger

	projectIDMaybe := r.URL.Query().Get("projectID")
	if projectIDMaybe == "" {
		logger.Error("missing project ID")
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	projectIDInst, err := projectID.New(projectIDMaybe)
	if err != nil {
		logger.Error("invalid project ID", "error", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	uow, err := appCtx.UowFactory.Begin()
	if err != nil {
		logger.Error("database access failed", "error", err)
		http.Error(w, "Failed to access database", http.StatusInternalServerError)
		return
	}
	defer uow.Rollback()

	project, err := appCtx.ProjectDB.GetByID(projectIDInst)
	if err != nil {
		logger.Error("project not found", "projectID", projectIDMaybe, "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if project == nil {
		logger.Error("project not found", "projectID", projectIDMaybe)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	data := Data{
		Project:     project.EnsureComputed(),
		ProjectPage: projectRoutes.ToProjectPage(projectIDInst),
	}

	page.Respond(static.GetSiblingPath("page.html"), data)(w, r)
}

func respondPost(appCtx *appContext.AppCtx, w http.ResponseWriter, r *http.Request) {
	req := reqCtx.FromHttpRequest(appCtx, r)
	logger := req.Logger

	logger.Info("handling project delete request")

	// Handle form submission
	if err := r.ParseForm(); err != nil {
		logger.Error("failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	projectIDMaybe := r.FormValue("projectID")
	if projectIDMaybe == "" {
		logger.Error("missing project ID")
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	projectIDInst, err := projectID.New(projectIDMaybe)
	if err != nil {
		logger.Error("invalid project ID", "error", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	confirmDelete := r.FormValue("confirmDelete")
	if confirmDelete != "DELETE" {
		logger.Error("delete confirmation not provided", "confirmation", confirmDelete)
		http.Error(w, "You must type DELETE to confirm", http.StatusBadRequest)
		return
	}

	// Get existing project
	uow, err := appCtx.UowFactory.Begin()
	if err != nil {
		logger.Error("failed to begin transaction", "error", err)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}
	defer uow.Rollback()

	existingProject, err := appCtx.ProjectDB.GetByID(projectIDInst)
	if err != nil {
		logger.Error("project not found", "projectID", projectIDMaybe, "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if existingProject == nil {
		logger.Error("project not found", "projectID", projectIDMaybe)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	logger.Info("deleting project", "projectID", projectIDInst)

	if err = appCtx.ProjectDB.ZapByID(uow, projectIDInst); err != nil {
		logger.Error("failed to delete project", "error", err)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	if err = uow.Commit(); err != nil {
		logger.Error("failed to commit transaction", "error", err)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	logger.Info("project deleted successfully", "projectID", projectIDInst)
	http.Redirect(w, r, projectRoutes.ToProjectListPage(), http.StatusSeeOther)
}
