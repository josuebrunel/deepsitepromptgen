package prompt

import (
	"bytes"
	"strings"
	"text/template"
)

type WebsiteData struct {
	Name         string   `json:"name" form:"name"`
	Description  string   `json:"description" form:"description"`
	CSSFramework string   `json:"cssFramework" form:"cssFramework"`
	Instructions string   `json:"instructions" form:"instructions"`
	Pages        []string `json:"pages" form:"pages"`
}

func Generate(data WebsiteData) (string, error) {
	tpl, err := template.New("website").Parse(websitePrompt)
	if err != nil {
		return "", err
	}

	// Filter empty pages
	var validPages []string
	for _, p := range data.Pages {
		if strings.TrimSpace(p) != "" {
			validPages = append(validPages, p)
		}
	}
	data.Pages = validPages

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const websitePrompt = `
You are a world-class UI/UX Designer and Frontend Architect known for building award-winning, accessible, and SEO-optimized web applications.
Your task is to generate a stunning, single-page website acting as a Single Page Application (SPA).

**Input Data:**
- **Application Name:** {{.Name}}
- **Application Description:** {{.Description}}
- **CSS Framework:** {{.CSSFramework}} (Use a CDN. If Tailwind, use the script tag for rapid prototyping. If Bootstrap, use version 5+).
- **Special Instructions:** {{.Instructions}}
- **Page Structure:**
{{- range .Pages }}
  - {{ . }}
{{- end }}

**Design & Aesthetic Guidelines:**
1.  **Visual Style:** Create a "High-End" look using generous whitespace, consistent padding, and modern typography.
2.  **Color Palette & Compliance:** Derive a professional palette based on {{.Name}}. **Crucial:** You must ensure **WCAG 2.1 AA compliance**. All text foreground/background combinations must meet a minimum contrast ratio of 4.5:1. Avoid subtle greys on white backgrounds if they fail this ratio.
3.  **Typography:** Import high-quality Google Fonts (Display font for headings, readable Sans-Serif for body).
4.  **Imagery:** Use high-quality placeholder images (Unsplash Source) relevant to the content.
5.  **Icons:** Include a CDN for an icon library (e.g., FontAwesome or Heroicons).
6. **Background Artistry (Crucial):** Do not use flat, solid backgrounds. Implement a sophisticated background strategy:
    *   **Mesh Gradients & Depth:** Use multi-layered radial gradients or CSS-animated organic "blobs" to create a sense of movement and depth.
    *   **Glassmorphism:** Use *backdrop-filter: blur(12px)* and semi-transparent backgrounds for cards and navigation to allow the background colors to bleed through elegantly.
    *   **Textures:** Incorporate subtle SVG patterns (e.g., topographic lines, architectural grids, or noise textures) at low opacity to add tactile quality.
    *   **Section Transitions:** Use angled clips or SVG wave dividers between sections for a seamless flow

**Technical Requirements:**
1.  **Single HTML File:** All CSS, JS, and HTML must be contained in one file.
2.  **SPA Navigation:**
    * Sticky or accessible navigation.
    * JS-based section hiding/showing (no page reloads).
    * Dynamic "active" state updates on nav links.
3.  **Responsiveness:** Mobile-first approach.
4.  **SEO & JSON-LD:**
    * Use semantic HTML5 tags (*<header>*, *<main>*, *<article>*, *<footer>*).
    * **JSON-LD:** Generate a *<script type="application/ld+json">* block in the *<head>*. Create a Schema.org object (e.g., *LocalBusiness*, *Organization*, or *Product*) that accurately describes **{{.Name}}** using the provided description.

**Accessibility (WCAG 2.1 AA) Strictness:**
1.  **Focus Indicators:** Do NOT remove default focus outlines unless replacing them with a custom, high-visibility focus state (e.g., a thick ring or border color change). All interactive elements (buttons, inputs, links) must be clearly visible when navigated via keyboard.
2.  **ARIA:** Use *aria-label* or *aria-labelledby* for any interactive element that lacks visible text (like icon-only buttons).
3.  **Semantic Structure:** Use proper heading hierarchy (*h1* -> *h2* -> *h3*).

**Functionality:**
1.  **Notification System:** Implement a non-intrusive "Toast" notification system (JS).
2.  **Interactive Elements:** Ensure smooth transitions (*transition: all 0.3s ease*).

**Content Strategy:**
1.  **No Lorem Ipsum:** Generate realistic, persuasive marketing copy relevant to **{{.Name}}**.

**Output:**
Provide **only** the complete HTML code block.
`
