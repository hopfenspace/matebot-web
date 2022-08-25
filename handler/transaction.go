package handler

import "github.com/labstack/echo/v4"

func (a *API) SendTransaction(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) ConsumeTransaction(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}
