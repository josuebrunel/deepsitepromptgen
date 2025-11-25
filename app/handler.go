package app

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/josuebrunel/gopkg/etr"
	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
)

const websitePrompt = `You are an expert web developer and designer.
Using the following inputs, generate a complete, responsive, SEO-friendly, and accessible HTML/CSS/JS website:

- **Application Name:** {{.Name}}
- **Application Description:** {{.Description}}
- **CSS Framework:** {{.CSSFramework}} (use CDN link if available)
- **Special Instructions:** {{.Instructions}}
- **Pages:**
{{- range .Pages }}
  - {{ . }}
{{- end }}

Requirements:
1. Include a clean navigation bar linking to all pages.
2. Use semantic HTML5 and ensure accessibility (ARIA labels, alt attributes, proper heading structure).
3. Apply the specified CSS framework for styling.
4. Create a consistent layout across all pages, with a header, main content area, and footer.
5. Add placeholder content relevant to the application description.
6. Optimize for SEO with meta tags, titles, and descriptions based on the application name and description.
7. Ensure mobile-first responsiveness.
8. If special instructions are given, integrate them thoughtfully into the design.
9. Produce a single html page
10. Emulate navigation between components with a show/hide javascript function
11. Add a notification system
`

type WebsiteData struct {
	Name         string   `json:"name" form:"name"`
	Description  string   `json:"description" form:"description"`
	CSSFramework string   `json:"css_framework" form:"css_framework"`
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
