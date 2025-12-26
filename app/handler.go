package app

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/josuebrunel/gopkg/etr"
	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
)

type WebsiteData struct {
	Name         string   `json:"name" form:"name"`
	Description  string   `json:"description" form:"description"`
	CSSFramework string   `json:"cssFramework" form:"cssFramework"`
	Instructions string   `json:"instructions" form:"instructions"`
	Pages        []string `json:"pages" form:"pages"`
}

func HandlerIndex() echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == http.MethodGet {
			return etr.Render(c, http.StatusOK, Index(), nil)
		}

		var req WebsiteData
		if err := c.Bind(&req); err != nil {
			xlog.Error("bind error", "error", err)
			return err
		}
		xlog.Debug("request data", "request", req)
		var buf bytes.Buffer
		tpl := template.Must(template.New("website").Parse(websitePrompt))
		if err := tpl.Execute(&buf, req); err != nil {
			xlog.Error("template error", "error", err)
			return err
		}
		return etr.Render(c, http.StatusOK, Prompt(req.Name, buf.String()), nil)
	}
}
