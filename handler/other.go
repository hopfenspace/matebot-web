package handler

import "github.com/labstack/echo/v4"

func (a *API) Consumables(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) Applications(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) Balance(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) Blame(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) Zwegat(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}
