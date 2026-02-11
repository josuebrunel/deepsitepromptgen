package prompt

import (
	"dppg/internal/db/repository"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/josuebrunel/ezauth"
	ezhandler "github.com/josuebrunel/ezauth/pkg/handler"
	"github.com/josuebrunel/gopkg/etr"
	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	Auth *ezauth.EzAuth
	Repo *repository.Repository
}

func NewHandler(auth *ezauth.EzAuth, repo *repository.Repository) *Handler {
	return &Handler{Auth: auth, Repo: repo}
}

func Mount(e *echo.Echo, auth *ezauth.EzAuth, repo *repository.Repository) {
	h := NewHandler(auth, repo)
	g := e.Group("/prompts",
		echo.WrapMiddleware(auth.Handler.Session.LoadAndSave),
		echo.WrapMiddleware(auth.LoadUserMiddleware),
		echo.WrapMiddleware(auth.LoginRequiredMiddleware),
	)

	g.AddRoute(echo.Route{Method: "GET", Path: "", Handler: h.List(), Name: "List"})
	g.AddRoute(echo.Route{Method: "POST", Path: "/new", Handler: h.Create(), Name: "Create"})
	g.AddRoute(echo.Route{Method: "GET", Path: "/new", Handler: h.Create(), Name: "Create"})
	g.AddRoute(echo.Route{Method: "GET", Path: "/:id", Handler: h.View(), Name: "View"})
	g.AddRoute(echo.Route{Method: "DELETE", Path: "/:id", Handler: h.Delete(), Name: "Delete"})
}

func (h *Handler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		// userID, err := h.Auth.GetUserID(c.Request().Context())
		// if err != nil || userID == "" {
		// 	return c.Redirect(http.StatusFound, "/auth/login")
		// }
		user, err := h.Auth.GetSessionUser(c.Request().Context())
		if err != nil || user == nil {
			xlog.Error("failed to get session's user", "error", err)
			return c.Redirect(http.StatusFound, "/")
		}
		uid, err := uuid.FromString(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
		}

		prompts, err := h.Repo.ListByUserID(c.Request().Context(), uid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return etr.Render(c, http.StatusOK, List(prompts), nil)
	}
}

func (h *Handler) Create() echo.HandlerFunc {
	return func(c *echo.Context) error {

		if c.Request().Method == http.MethodGet {
			return etr.Render(c, http.StatusOK, Form(), nil)
		}

		userID, err := h.Auth.GetUserID(c.Request().Context())
		if err != nil || userID == "" {
			xlog.Error("failed to get user id", "error", err)
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		uid, err := uuid.FromString(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
		}

		var req WebsiteData
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		promptContent, err := Generate(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		prompt, err := h.Repo.Create(c.Request().Context(), uid, req.Name, promptContent)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.Redirect(http.StatusFound, "/prompts/"+prompt.ID.String())
	}
}

func (h *Handler) View() echo.HandlerFunc {
	return func(c *echo.Context) error {
		userID, err := ezhandler.GetUserID(c.Request().Context())
		if err != nil || userID == "" {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		uid, err := uuid.FromString(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
		}

		idParam := c.Param("id")
		id, err := uuid.FromString(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid prompt ID")
		}

		prompt, err := h.Repo.GetByID(c.Request().Context(), id, uid)
		if err != nil {
			return etr.Render(c, http.StatusNotFound, NotFound(), nil)
		}

		return etr.Render(c, http.StatusOK, Prompt(prompt.Name, prompt.Content), nil)
	}
}

func (h *Handler) Delete() echo.HandlerFunc {
	return func(c *echo.Context) error {
		userID, err := ezhandler.GetUserID(c.Request().Context())
		if err != nil || userID == "" {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		uid, err := uuid.FromString(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID")
		}

		idParam := c.Param("id")
		id, err := uuid.FromString(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid prompt ID")
		}

		if err := h.Repo.Delete(c.Request().Context(), id, uid); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}
