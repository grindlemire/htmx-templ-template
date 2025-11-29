package router

import (
	"context"
	"net/http"

	"github.com/grindlemire/gothem-stack/pkg/auth"
	"github.com/grindlemire/gothem-stack/pkg/handler"
	"github.com/grindlemire/gothem-stack/web"

	"github.com/grindlemire/graft"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

const ID = graft.ID("router")

type Output struct {
	http.Handler
}

func init() {
	graft.Register(graft.Node[Output]{
		ID:        ID,
		Cacheable: true,
		Run:       run,
	})
}

func run(ctx context.Context) (Output, error) {
	e := echo.New()

	e.Use(
		// recover from panics and create errors from them
		middleware.Recover(),
		// other global middleware goes here
	)

	// register the customer pages and components
	homeHandler, err := handler.NewHomeHandler()
	if err != nil {
		return Output{}, err
	}
	homeHandler.RegisterRoutes(
		e.Group("", auth.Middleware()),
	)

	// register the static assets like the favicon and the css
	err = web.RegisterStaticAssets(e)
	if err != nil {
		return Output{}, err
	}

	// all other routes should return not found. This should be the last registered route in the list
	e.HTTPErrorHandler = handler.Error
	e.Add(echo.RouteNotFound, "/*", echo.HandlerFunc(func(c echo.Context) error {
		return echo.ErrNotFound.SetInternal(errors.Errorf("not found | uri=[%s]", c.Request().RequestURI))
	}), []echo.MiddlewareFunc{}...)

	return Output{
		Handler: e.Server.Handler,
	}, nil
}
