package app

import (
	"dppg/internal/db/repository"
	"net/http"

	"github.com/josuebrunel/ezauth"
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

func (h *Handler) Index() echo.HandlerFunc {
	return func(c *echo.Context) error {
		if authErrMsg := h.Auth.GetErrorMessage(c.Request().Context()); authErrMsg != "" {
			xlog.Error("auth error message", "message", authErrMsg)
			return etr.Render(c, http.StatusOK, Index(false), nil)
		}
		return etr.Render(c, http.StatusOK, Index(false), nil)
	}
}

func (h *Handler) SignUp() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return etr.Render(c, http.StatusOK, Index(true), nil)
	}
}
