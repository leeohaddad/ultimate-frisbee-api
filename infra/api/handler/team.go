package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	handlerParam "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/param"
	handlerResult "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/result"

	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/payload"

	repositoryPort "github.com/leeohaddad/ultimate-frisbee-api/domain/port/repository"
	domainService "github.com/leeohaddad/ultimate-frisbee-api/domain/service"

	"github.com/labstack/echo/v4"
)

// GetAllTeamsEchoHandlerV1 is the adapter from the Echo ecosystem to the GetAllTeams handler.
func GetAllTeamsEchoHandlerV1(param handlerParam.GetAllTeamsHandlerV1) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()

		return DispatchEchoResponseFromHandlerResult(echoContext, GetAllTeamsHandlerV1(requestContext, param).HTTP)
	}
}

// GetAllTeamsHandlerV1 is the entry point to the application's logic for fetching a list of existing teams.
func GetAllTeamsHandlerV1(
	context context.Context,
	param handlerParam.GetAllTeamsHandlerV1,
) handlerResult.GetAllTeamsHandlerV1 {
	result, err := domainService.GetAllTeams(context, domainServiceParam.GetAllTeams{
		Repository: param.Repository,
	})
	if err != nil {
		return handlerResult.GetAllTeamsHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("failed to get all teams from domain service: %s", err.Error()),
			},
		}
	}

	return handlerResult.GetAllTeamsHandlerV1{
		HTTP: handlerResult.HTTP{
			StatusCode:   http.StatusOK,
			ResponseType: handlerResult.ResponseBodyTypes.JSON,
			JSONResponse: payload.TeamEntitiesToTeams(result.Teams),
		},
	}
}

// GetTeamByNameEchoHandlerV1 is the adapter from the Echo ecosystem to the GetTeamByName handler.
func GetTeamByNameEchoHandlerV1(param handlerParam.GetTeamByNameHandlerV1) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()
		param.Name = echoContext.Param("name")

		return DispatchEchoResponseFromHandlerResult(echoContext, GetTeamByNameHandlerV1(requestContext, param).HTTP)
	}
}

// GetTeamByNameHandlerV1 is the entry point to the application's logic of fetching an specific team by its name.
func GetTeamByNameHandlerV1(
	context context.Context,
	param handlerParam.GetTeamByNameHandlerV1,
) handlerResult.GetTeamByNameHandlerV1 {
	result, err := domainService.GetTeamByName(context, domainServiceParam.GetTeamByName{
		Name:       param.Name,
		Repository: param.Repository,
	})
	if err != nil {
		return handlerResult.GetTeamByNameHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("failed to search team by name '%s' from domain service: %s", param.Name, err.Error()),
			},
		}
	}

	if result.Team == nil {
		return handlerResult.GetTeamByNameHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusNotFound,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("no team with name '%s' was found in the repository", param.Name),
			},
		}
	}

	return handlerResult.GetTeamByNameHandlerV1{
		HTTP: handlerResult.HTTP{
			StatusCode:   http.StatusOK,
			ResponseType: handlerResult.ResponseBodyTypes.JSON,
			JSONResponse: payload.TeamEntityToTeam(result.Team),
		},
	}
}

// CreateTeamEchoHandlerV1 is the adapter from the Echo ecosystem to the CreateTeam handler.
func CreateTeamEchoHandlerV1(param handlerParam.CreateTeamHandlerV1) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()

		var team payload.Team
		err := echoContext.Bind(&team)
		if err != nil {
			return DispatchEchoResponseFromString(echoContext, http.StatusBadRequest, "invalid payload format")
		}
		param.Payload = team

		return DispatchEchoResponseFromHandlerResult(echoContext, CreateTeamHandlerV1(requestContext, param).HTTP)
	}
}

// CreateTeamHandlerV1 is the entry point to the application's logic of creating a new team.
func CreateTeamHandlerV1(context context.Context, param handlerParam.CreateTeamHandlerV1) handlerResult.CreateTeamHandlerV1 {
	paramsAreValid, invalidParamsMessage := payload.ValidateCreateTeamInput(&param.Payload)
	if !paramsAreValid {
		return handlerResult.CreateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusBadRequest,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: invalidParamsMessage,
			},
		}
	}

	result, err := domainService.CreateTeam(context, domainServiceParam.CreateTeam{
		Team:       payload.TeamToTeamEntity(param.Payload),
		Repository: param.Repository,
	})
	if err != nil {
		// map repository sentinel error to HTTP 409 Conflict
		if errors.Is(err, repositoryPort.ErrAlreadyExists) {
			return handlerResult.CreateTeamHandlerV1{
				HTTP: handlerResult.HTTP{
					StatusCode:     http.StatusConflict,
					ResponseType:   handlerResult.ResponseBodyTypes.String,
					StringResponse: fmt.Sprintf("team with name '%s' already exists", param.Payload.Name),
				},
			}
		}

		return handlerResult.CreateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("failed to create team with name '%s' in domain service: %s", param.Payload.Name, err.Error()),
			},
		}
	}
	if result.Team == nil {
		return handlerResult.CreateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("no team with name '%s' was created", param.Payload.Name),
			},
		}
	}
	return handlerResult.CreateTeamHandlerV1{
		HTTP: handlerResult.HTTP{
			StatusCode:   http.StatusCreated,
			ResponseType: handlerResult.ResponseBodyTypes.JSON,
			JSONResponse: payload.TeamEntityToTeam(result.Team),
		},
	}
}

// UpdateTeamEchoHandlerV1 is the adapter from the Echo ecosystem to the UpdateTeam handler.
func UpdateTeamEchoHandlerV1(param handlerParam.UpdateTeamHandlerV1) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()
		param.Name = echoContext.Param("name")

		var team payload.Team
		err := echoContext.Bind(&team)
		if err != nil {
			return DispatchEchoResponseFromString(echoContext, http.StatusBadRequest, "invalid payload format")
		}
		param.Payload = team

		return DispatchEchoResponseFromHandlerResult(echoContext, UpdateTeamHandlerV1(requestContext, param).HTTP)
	}
}

// UpdateTeamHandlerV1 is the entry point to the application's logic of updating info of an existing team.
func UpdateTeamHandlerV1(context context.Context, param handlerParam.UpdateTeamHandlerV1) handlerResult.UpdateTeamHandlerV1 {
	paramsAreValid, invalidParamsMessage := payload.ValidateUpdateTeamInput(&param.Payload, param.Name)
	if !paramsAreValid {
		return handlerResult.UpdateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusBadRequest,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: invalidParamsMessage,
			},
		}
	}
	param.Payload.Name = param.Name

	result, err := domainService.UpdateTeam(context, domainServiceParam.UpdateTeam{
		Team:              payload.TeamToTeamEntity(param.Payload),
		UpdatedAttributes: payload.GetFilledTeamAttributesForUpdate(&param.Payload),
		Repository:        param.Repository,
	})
	if err != nil {
		return handlerResult.UpdateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("failed to update team with name '%s' in domain service: %s", param.Name, err.Error()),
			},
		}
	}

	if result.Team == nil {
		return handlerResult.UpdateTeamHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusNotFound,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("no team with name '%s' was found", param.Name),
			},
		}
	}
	return handlerResult.UpdateTeamHandlerV1{
		HTTP: handlerResult.HTTP{
			StatusCode:   http.StatusOK,
			ResponseType: handlerResult.ResponseBodyTypes.JSON,
			JSONResponse: payload.TeamEntityToTeam(result.Team),
		},
	}
}
