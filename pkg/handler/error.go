package handler

import (
	"encoding/json"
	"net/http"

	"github.com/grindlemire/gothem-stack/web/pages/authfail"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Error(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if ok {
		// If there is an internal error then use that for printing
		if he.Internal != nil {
			err = he.Internal
		}
	} else {
		he = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: err,
		}
	}

	code := he.Code
	message := he.Message

	// only log unauthorized errors at the debug level
	if he.Code != http.StatusUnauthorized {
		zap.S().Error(err)
	} else {
		zap.S().Debug(err)
		// render the auth fail page for unauthorized errors
		c.Response().WriteHeader(http.StatusUnauthorized)
		renderErr := render(c, authfail.Page())
		if renderErr != nil {
			zap.S().Error(errors.Wrap(renderErr, "rendering auth fail page"))
		}
		return
	}

	switch m := he.Message.(type) {
	case string:
		message = echo.Map{"message": m}
	case json.Marshaler:
		// do nothing - this type knows how to format itself to JSON
	case error:
		message = echo.Map{"message": m.Error()}
	}

	// Send response
	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
		if err != nil {
			zap.S().Error(errors.Wrap(err, "sending no content"))
		}
		return
	}

	err = c.JSON(code, message)
	if err != nil {
		zap.S().Error(errors.Wrap(err, "marshalling json payload"))
	}
}
