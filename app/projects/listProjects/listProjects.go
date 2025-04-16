package listProjects

import (
	"imageresizerservice/app/ctx/appCtx"
	"imageresizerservice/app/ctx/reqCtx"
	"imageresizerservice/app/home/homeRoutes"
	"imageresizerservice/app/projects/project"
	"imageresizerservice/app/projects/projectRoutes"
	"imageresizerservice/app/ui/page"
	"imageresizerservice/app/users/userAccount/userAccountRoutes"
	"imageresizerservice/library/static"
	"net/http"
)

func Router(mux *http.ServeMux, ac *appCtx.AppCtx) {
	mux.HandleFunc(projectRoutes.ProjectListPage, Respond(ac))
}

type Data struct {
	BackURL           string
	GetUserAccountURL string
	Projects          []*project.Project
	CreateURL         string
}

func Respond(ac *appCtx.AppCtx) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := reqCtx.FromHttpRequest(ac, r)
		logger := req.Logger
		createdByUserID := req.UserSession.UserID

		logger.Info("projectListPage", "userID", createdByUserID)

		uow, err := ac.UowFactory.Begin()
		if err != nil {
			logger.Error("database access failed", "error", err)
			http.Error(w, "Failed to access database", http.StatusInternalServerError)
			return
		}
		defer uow.Rollback()

		projects, err := ac.ProjectDB.GetByCreatedByUserID(createdByUserID)
		if err != nil {
			logger.Error("failed to fetch projects", "userID", createdByUserID, "error", err)
			http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
			return
		}

		for _, project := range projects {
			project.EnsureComputed()
		}

		data := Data{
			BackURL:           homeRoutes.HomePage,
			GetUserAccountURL: userAccountRoutes.UserAccountPage,
			Projects:          projects,
			CreateURL:         projectRoutes.ToProjectCreate(),
		}

		logger.Info("rendering project list page", "projectCount", len(projects))
		page.Respond(static.GetSiblingPath("page.html"), data)(w, r)
	}
}
