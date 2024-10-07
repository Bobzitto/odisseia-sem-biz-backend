package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) WriteJSON(c echo.Context, status int, data interface{}, headers ...http.Header) error {
	//set content type
	c.Response().Header().Set("Content-Type", "application/json")

	//add extra headers if necessary

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Response().Header()[key] = value
		}
	}

	//set status code
	c.Response().WriteHeader(status)

	//encode
	return json.NewEncoder(c.Response()).Encode(data)
}

func (app *application) ReadJSON(c echo.Context, dst interface{}) error {
	//limit size of request
	maxBytes := 1048576
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, int64(maxBytes))
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(dst)
	if err != nil {
		return err
	}
	if decoder.More() {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

func (app *application) errorJSON(c echo.Context, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.WriteJSON(c, statusCode, payload)
}
