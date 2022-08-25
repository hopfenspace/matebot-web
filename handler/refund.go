package handler

import "github.com/labstack/echo/v4"

func (a *API) NewRefund(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) ApproveRefund(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) DisapproveRefund(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) AbortRefund(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) OpenRefunds(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}

func (a *API) AllRefunds(c echo.Context) error {
	return c.JSON(501, GenericResponse{"Not implemented yet."})
}
