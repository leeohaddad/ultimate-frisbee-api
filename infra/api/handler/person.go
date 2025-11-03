package handler

import (
	"context"
	"fmt"
	"net/http"

	handlerParam "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/param"
	handlerResult "github.com/leeohaddad/ultimate-frisbee-api/infra/api/handler/result"

	domainServiceParam "github.com/leeohaddad/ultimate-frisbee-api/domain/service/param"

	"github.com/leeohaddad/ultimate-frisbee-api/infra/api/payload"

	domainService "github.com/leeohaddad/ultimate-frisbee-api/domain/service"

	"github.com/labstack/echo/v4"
)

// GetAllPeopleEchoHandlerV1 is the adapter from the Echo ecosystem to the GetAllPeople handler.
func GetAllPeopleEchoHandlerV1(param handlerParam.GetAllPeopleHandlerV1) echo.HandlerFunc {
	return func(echoContext echo.Context) error {
		requestContext := echoContext.Request().Context()

		return DispatchEchoResponseFromHandlerResult(echoContext, GetAllPeopleHandlerV1(requestContext, param).HTTP)
	}
}

// GetAllPeopleHandlerV1 is the entry point to the application's logic for fetching a list of existing People.
func GetAllPeopleHandlerV1(
	context context.Context,
	param handlerParam.GetAllPeopleHandlerV1,
) handlerResult.GetAllPeopleHandlerV1 {
	result, err := domainService.GetAllPeople(context, domainServiceParam.GetAllPeople{
		Repository: param.Repository,
	})
	if err != nil {
		return handlerResult.GetAllPeopleHandlerV1{
			HTTP: handlerResult.HTTP{
				StatusCode:     http.StatusInternalServerError,
				ResponseType:   handlerResult.ResponseBodyTypes.String,
				StringResponse: fmt.Sprintf("failed to get all People from domain service: %s", err.Error()),
			},
		}
	}

	return handlerResult.GetAllPeopleHandlerV1{
		HTTP: handlerResult.HTTP{
			StatusCode:   http.StatusOK,
			ResponseType: handlerResult.ResponseBodyTypes.JSON,
			JSONResponse: payload.PersonEntitiesToPeople(result.People),
		},
	}
}
