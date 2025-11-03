package api

import (
	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler"
	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/param"
)

func (app *App) configureRoutes() {
	app.configureRoutesV1()
}

func (app *App) configureRoutesV1() {
	v1RouterGroup := app.router.Group("/v1")

	v1RouterGroup.GET("/health/", handler.HealthCheckHandlerV1())

	// Teams
	v1RouterGroup.GET("/teams/", handler.GetAllTeamsEchoHandlerV1(
		param.GetAllTeamsHandlerV1{
			Repository: app.repositories.Team,
		},
	))
	v1RouterGroup.GET("/teams/:name/", handler.GetTeamByNameEchoHandlerV1(
		param.GetTeamByNameHandlerV1{
			Repository: app.repositories.Team,
		},
	))
	v1RouterGroup.POST("/teams/", handler.CreateTeamEchoHandlerV1(
		param.CreateTeamHandlerV1{
			Repository: app.repositories.Team,
		},
	))
	v1RouterGroup.PUT("/teams/:name/", handler.UpdateTeamEchoHandlerV1(
		param.UpdateTeamHandlerV1{
			Repository: app.repositories.Team,
		},
	))

	// People
	v1RouterGroup.GET("/people/", handler.GetAllPeopleEchoHandlerV1(
		param.GetAllPeopleHandlerV1{
			Repository: app.repositories.Person,
		},
	))
}
