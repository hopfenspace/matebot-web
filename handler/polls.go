package handler

import "github.com/labstack/echo/v4"

func (a *API) NewPoll(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) ApprovePoll(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) DisapprovePoll(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) AbortPoll(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) OpenPolls(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) AllPolls(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}
