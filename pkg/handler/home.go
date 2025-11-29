package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/grindlemire/gothem-stack/web/pages/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	// you could put a database handle here or any dependencies you want
}

func NewHomeHandler() (h *HomeHandler, err error) {
	return &HomeHandler{}, nil
}

// RegisterRoutes registers all the subroutes for the home handler to manage
func (h *HomeHandler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.RenderHomepage)
	g.GET("/random-string", h.GetRandomString)
}

func (h *HomeHandler) RenderHomepage(c echo.Context) error {
	return render(c, home.Page(c.Request().URL.Query().Get("username")))
}

func (h *HomeHandler) GetRandomString(c echo.Context) error {
	time.Sleep(750 * time.Millisecond)
	return render(c, home.RandomString(uuid.NewString()))
}
