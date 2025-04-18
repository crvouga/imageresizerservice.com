package apiDocsPage

import (
	"imageresizerservice/app/ctx/appCtx"
	"imageresizerservice/app/ctx/reqCtx"
	"imageresizerservice/app/projects/project"
	"imageresizerservice/app/projects/project/projectID"
	"imageresizerservice/app/projects/projectRoutes"
	"imageresizerservice/app/users/userID"
	"net/http"
)

type ProjectSelect struct {
	Projects         []*project.Project
	CreateProjectURL string
	ProjectID        projectID.ProjectID
}

func GetProjectSelect(ac *appCtx.AppCtx, r *http.Request) ProjectSelect {
	rc := reqCtx.FromHttpRequest(ac, r)
	return ProjectSelect{
		Projects:         getProjects(ac, rc.UserSession.UserID),
		ProjectID:        projectID.ProjectID(r.URL.Query().Get("projectID")),
		CreateProjectURL: projectRoutes.ToCreateProject(),
	}
}

func getProjects(ac *appCtx.AppCtx, userID userID.UserID) []*project.Project {
	projects, err := ac.ProjectDB.GetByCreatedByUserID(userID)

	if err != nil {
		return []*project.Project{}
	}

	return projects
}
